package database

import (
	"database/sql"
	"devices/internal/types"
	"fmt"
)

type DB struct {
	db *sql.DB
}

func New(db *sql.DB) DB {
	return DB{db: db}
}

func (d DB) GetDevices() ([]types.Device, error) {
	var devices []types.Device

	rows, err := d.db.Query("SELECT * FROM devices")
	if err != nil {
		return nil, fmt.Errorf("GetDevices: %v", err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var device types.Device
		if err := rows.Scan(&device.ID, &device.Name, &device.IsOccupied); err != nil {
			return nil, fmt.Errorf("GetDevices: %v", err)
		}
		devices = append(devices, device)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetDevices: %v", err)
	}
	return devices, nil
}

func (d DB) GetDevice(deviceID string) (types.Device, error) {
	dev := types.Device{}

	row := d.db.QueryRow("SELECT id, email FROM users WHERE id=$1;", deviceID)
	err := row.Scan(dev.ID, dev.Name, dev.Name)
	if err != nil {
		return types.Device{}, err
	}

	return dev, nil
}
