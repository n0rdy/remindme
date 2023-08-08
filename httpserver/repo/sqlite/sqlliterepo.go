package sqlite

import (
	"database/sql"
	"errors"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpserver/repo"
	"n0rdy.me/remindme/logger"
	"n0rdy.me/remindme/utils"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteReminderRepo struct {
	db *sql.DB
}

func NewSqliteReminderRepo() (repo.ReminderRepo, error) {
	db, err := sql.Open("sqlite3", utils.GetOsSpecificLogsDir()+"remindme.db")
	if err != nil {
		return nil, err
	}

	logger.Log("SQLite DB started")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS reminders (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    message TEXT NOT NULL,
		    remind_at INTEGER NOT NULL
		);
	`)

	if err != nil {
		return nil, err
	}

	logger.Log("SQLite DB schema created")

	return &sqliteReminderRepo{db: db}, nil
}

func (repo *sqliteReminderRepo) Add(reminder common.Reminder) error {
	_, err := repo.db.Exec(`
		INSERT INTO reminders (message, remind_at) VALUES (?, ?);
	`, reminder.Message, reminder.RemindAt.Unix())
	return err
}

func (repo *sqliteReminderRepo) Update(reminder common.Reminder) error {
	_, err := repo.db.Exec(`
		UPDATE reminders SET message = ?, remind_at = ? WHERE id = ?;
	`, reminder.Message, reminder.RemindAt.Unix(), reminder.ID)
	return err
}

func (repo *sqliteReminderRepo) List() ([]common.Reminder, error) {
	rows, err := repo.db.Query(`
		SELECT id, message, remind_at FROM reminders;
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reminders := make([]common.Reminder, 0)
	for rows.Next() {
		var id int
		var message string
		var remindAt int64

		err := rows.Scan(&id, &message, &remindAt)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, common.Reminder{
			ID:       id,
			Message:  message,
			RemindAt: time.Unix(remindAt, 0),
		})
	}
	return reminders, nil
}

func (repo *sqliteReminderRepo) Get(id int) (*common.Reminder, error) {
	row := repo.db.QueryRow(`
		SELECT id, message, remind_at FROM reminders WHERE id = ?;
	`, id)

	var reminderId int
	var reminderMessage string
	var remindAtUnix int64

	err := row.Scan(&reminderId, &reminderMessage, &remindAtUnix)
	if err != nil {
		// no rows required a special handling as it's not an error, but rather a DB state
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &common.Reminder{
		ID:       reminderId,
		Message:  reminderMessage,
		RemindAt: time.Unix(remindAtUnix, 0),
	}, nil
}

func (repo *sqliteReminderRepo) DeleteAll() error {
	_, err := repo.db.Exec(`
		DELETE FROM reminders;
	`)
	return err
}

func (repo *sqliteReminderRepo) Delete(id int) error {
	_, err := repo.db.Exec(`
		DELETE FROM reminders WHERE id = ?;
	`, id)
	return err
}

func (repo *sqliteReminderRepo) Exists(id int) (bool, error) {
	row := repo.db.QueryRow(`
		SELECT id FROM reminders WHERE id = ?;
	`, id)

	var foundId int
	err := row.Scan(&foundId)
	if err != nil {
		// no rows required a special handling as it's not an error, but rather a DB state
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (repo *sqliteReminderRepo) DeleteAllWithRemindAtBefore(threshold time.Time) ([]int, error) {
	rows, err := repo.db.Query(`
		SELECT id FROM reminders WHERE remind_at < ?;
	`, threshold.Unix())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int, 0)
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	_, err = repo.db.Exec(`
		DELETE FROM reminders WHERE remind_at < ?;
	`, threshold.Unix())
	return ids, err
}

func (repo *sqliteReminderRepo) GetRemindersAfter(threshold time.Time) ([]common.Reminder, error) {
	rows, err := repo.db.Query(`
		SELECT id, message, remind_at FROM reminders WHERE remind_at > ?;
	`, threshold.Unix())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reminders := make([]common.Reminder, 0)
	for rows.Next() {
		var id int
		var message string
		var remindAt int64

		err := rows.Scan(&id, &message, &remindAt)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, common.Reminder{
			ID:       id,
			Message:  message,
			RemindAt: time.Unix(remindAt, 0),
		})
	}
	return reminders, nil
}

func (repo *sqliteReminderRepo) Close() error {
	return repo.db.Close()
}
