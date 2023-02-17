//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.1.0
//     //             ///   ///  ///    Date        : 2022/02/18
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description :
//   ////                ///  ///
//     ////     //////////   ///
//        ////            ////
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package main

import (
	"time"
)

func GetINIString(section string, key string, defaultValue string) string {
	for _, s := range cfg.SectionStrings() {
		switch {
		case s != section:
			continue
		case cfg.Section(section).HasKey(key):
			return cfg.Section(section).Key(key).MustString(defaultValue)
		default:
			return defaultValue
		}
	}
	return defaultValue
}

func GetINIInt(section string, key string, defaultValue int) int {
	for _, s := range cfg.SectionStrings() {
		switch {
		case s != section:
			continue
		case cfg.Section(section).HasKey(key):
			return cfg.Section(section).Key(key).MustInt(defaultValue)
		default:
			return defaultValue
		}
	}
	return defaultValue
}

func GetINIDuration(section string, key string, defaultValue time.Duration) time.Duration {
	for _, s := range cfg.SectionStrings() {
		switch {
		case s != section:
			continue
		case cfg.Section(section).HasKey(key):
			return cfg.Section(section).Key(key).MustDuration(defaultValue)
		default:
			return defaultValue
		}
	}
	return defaultValue
}

func GetINIBoolean(section string, key string, defaultValue bool) bool {
	for _, s := range cfg.SectionStrings() {
		switch {
		case s != section:
			continue
		case cfg.Section(section).HasKey(key):
			return cfg.Section(section).Key(key).MustBool(defaultValue)
		default:
			return defaultValue
		}
	}
	return defaultValue
}

func CheckCredentials(section string, email string, password string) bool {
	switch {
	case GetINIString(section, "email", "") != email:
		return false
	case GetINIString(section, "password", "") == password:
		return true
	case GetINIString(section, "password", "") == Hash(password):
		return true
	default:
		return false
	}
}