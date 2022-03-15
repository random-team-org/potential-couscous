package server

import (
	"bytes"
	"database/sql"
	"devices/internal/types"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"devices/internal/server/database"

	"github.com/bradfitz/gomemcache/memcache"
)

type allocate struct {
	deviceID string
	userID   string
}

type Service struct {
	*sync.Mutex
	list map[string]*allocate // allocated devices
	db   *sql.DB
	mch  string // memcached host string
}

func NewService(dbopts, mch string) *Service {
	if mch == "" {
		mch = "127.0.0.1:11211"
	}
	service := Service{
		list: make(map[string]*allocate),
		db:   connectDB(dbopts),
		mch:  mch,
	}
	return &service
}
func (s *Service) GetDevices() ([]types.Device, error) {
	db := database.New(s.db)

	return db.GetDevices()
}

func (s *Service) AllocateDevice(deviceID, userID string) error {
	var device types.Device

	s.Lock()

	mc := memcache.New(s.mch)
	cacheKey := fmt.Sprintf("devices-%s-%s", deviceID, userID)
	it, err := mc.Get(cacheKey)
	if err == nil {
		if it != nil {
			_ = json.NewDecoder(bytes.NewBuffer(it.Value)).Decode(&device)
		}
	} else {
		device, err := database.New(s.db).GetDevice(deviceID)
		if err != nil {
			return err
		}

		b := bytes.NewBuffer(nil)
		err = json.NewEncoder(b).Encode(device)
		if err == nil {
			_ = mc.Set(&memcache.Item{
				Key:   cacheKey,
				Value: b.Bytes(),
			})
		}
		return err
	}

	_, ok := s.list[userID]
	if !ok {
		return errors.New("allocated")
	}

	s.list[userID] = &allocate{
		deviceID: deviceID,
		userID:   userID,
	}

	s.Unlock()
	return nil
}

func (s *Service) FreeDevice(deviceID string) error {
	s.Lock()

	s.list[deviceID] = nil

	s.Unlock()
	return nil
}

func connectDB(dbopts string) *sql.DB {
	db, err := sql.Open("mysql", dbopts)
	if err != nil {
		panic("error connecting to the database: " + err.Error())
	}

	pingErr := db.Ping()
	if pingErr != nil {
		panic("")
	}

	fmt.Println("Connected!")
	return db
}
