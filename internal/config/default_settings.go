package config

type defaultSettingKey uint

const (
	DATETIMEFORMAT defaultSettingKey = 0x0
)

var defaultSettings = map[defaultSettingKey]interface{}{
	DATETIMEFORMAT: "2006/01/02 15:04:05.999999999",
}
