package elasticsearch

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	xov1alpha1 "github.com/90poe/elasticsearch-objects-operator/api/v1alpha1"
	"github.com/90poe/elasticsearch-objects-operator/internal/consts"
)

func diffMappings(old map[string]interface{}, new map[string]interface{}) (bool, error) {
	newSetKeys := getKeysFromSettings("", new)
	for _, newSetKey := range newSetKeys {
		// Normalize all settings as strings
		oldVal, ok := getStringValueFromSettings(old, newSetKey)
		if !ok {
			// No such value in K8S settings
			return true, nil
		}
		newVal, ok := getStringValueFromSettings(new, newSetKey)
		if !ok {
			return false, fmt.Errorf("can't get new settings for path %s", newSetKey)
		}
		if oldVal != newVal {
			return true, nil
		}
	}
	return false, nil
}

func diffSettings(k8sSett *xov1alpha1.ESIndexSettings,
	servSettings map[string]interface{}, index bool) (bool, error) {
	k8sSettInt := Settings{
		Index: *k8sSett,
	}
	k8sSettJSON, err := json.Marshal(k8sSettInt)
	if err != nil {
		return false, fmt.Errorf("can't make current index settings JSON: %w", err)
	}
	var k8sSettMap map[string]interface{}
	err = json.Unmarshal(k8sSettJSON, &k8sSettMap)
	if err != nil {
		return false, fmt.Errorf("can't make modified index settings from JSON: %w", err)
	}
	servSetKeys := getKeysFromSettings("", servSettings)
	for _, servSetKey := range servSetKeys {
		// Normalize all settings as strings
		k8sVal, ok := getStringValueFromSettings(k8sSettMap, servSetKey)
		if !ok {
			// No such value in K8S settings
			continue
		}
		servVal, ok := getStringValueFromSettings(servSettings, servSetKey)
		if !ok {
			return false, fmt.Errorf("can't get server settings for path %s", servSetKey)
		}
		if k8sVal != servVal {
			if !index {
				// Template can have any settings
				return true, nil
			}
			if _, ok = consts.ESStaticSettings[servSetKey]; ok {
				return false, fmt.Errorf("can't change static setting %s from '%s' to '%s'",
					servSetKey, servVal, k8sVal)
			}
			return true, nil
		}
	}
	return false, nil
}

func addManagedBy2Interface(src string) (map[string]interface{}, error) {
	var inter map[string]interface{}
	err := json.Unmarshal([]byte(src), &inter)
	if err != nil {
		return nil, fmt.Errorf("can't json unmarshal mappings: %w", err)
	}
	meta, ok := inter["_meta"]
	if !ok {
		// Adding required managed-by
		inter["_meta"] = map[string]interface{}{
			"managed-by": "elasticsearch-objects-operator.xo.90poe.io",
		}
		return inter, nil
	}
	metaMap, ok := meta.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid _meta map")
	}
	metaMap["managed-by"] = "elasticsearch-objects-operator.xo.90poe.io"

	return inter, nil
}

// Function would suite both int32 and int64
func getInt64ValueFromSettings(settings map[string]interface{}, path string) (int64, bool) {
	val, ok := getValueFromSettings(settings, path)
	if !ok {
		return 0, ok
	}
	// First check - value IS int32
	if val, ok := val.(int); ok {
		return int64(val), ok
	}
	// Second check - value is String ?
	if val, ok := val.(string); ok {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return 0, false
		}
		return int64(intVal), true
	}
	return 0, false
}

func getStringValueFromSettings(settings map[string]interface{}, path string) (string, bool) {
	val, ok := getValueFromSettings(settings, path)
	if !ok {
		return "", ok
	}
	valStr, ok := val.(string)
	if !ok {
		// value is not String ? Make it string
		return fmt.Sprintf("%#v", val), true
	}
	return valStr, ok
}

func getValueFromSettings(settings map[string]interface{}, path string) (interface{}, bool) {
	points := strings.Split(path, ".")
	val, ok := settings[points[0]]
	if !ok {
		return nil, ok
	}
	if len(points) == 1 {
		// Last point
		return val, ok
	}
	nexSettings, ok := val.(map[string]interface{})
	if !ok {
		return nil, ok
	}
	return getValueFromSettings(nexSettings, strings.Join(points[1:], "."))
}

func getKeysFromSettings(prefix string, settings map[string]interface{}) []string {
	ret := []string{}
	for key, val := range settings {
		keyPath := key
		if len(prefix) != 0 {
			keyPath = strings.Join([]string{prefix, key}, ".")
		}
		if valMap, ok := val.(map[string]interface{}); ok {
			newKeys := getKeysFromSettings(keyPath, valMap)
			ret = append(ret, newKeys...)
			continue
		}
		ret = append(ret, keyPath)
	}
	return ret
}

func isManagedByESOperator(settings map[string]interface{}) bool {
	managedBy, ok := getStringValueFromSettings(settings, "_meta.managed-by")
	if !ok {
		return false
	}
	if managedBy != consts.ESManagedByValue {
		return false
	}
	return true
}
