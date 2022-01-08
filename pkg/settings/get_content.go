package settings

import (
	"encoding/json"
	"ferry/global/orm"
	"ferry/models/system"
)

func GetContent(classify int) (content map[string]interface{}, err error) {
	var (
		settings system.Settings
	)

	err = orm.Eloquent.Where("classify = ?", classify).Find(&settings).Error
	if err != nil {
		return
	}

	err = json.Unmarshal(settings.Content, &content)
	if err != nil {
		return
	}

	return
}

func GetContentByKey(classify int, key string) (value interface{}, err error) {
	var (
		content map[string]interface{}
	)

	content, err = GetContent(classify)
	if err != nil {
		return
	}

	value = content[key]

	return
}
