package virtualbox

type KeyValueMap map[string]string
type KeyValuesMap map[string]KeyValueMap

type KeyValues []KeyValue
type KeyValue struct {
	Key   string
	Value string
}

func (kvs *KeyValues) decodeBridgeIfs() (KeyValuesMap, bool) {

	ok := false
	kvm := make(KeyValuesMap)
	firstNic := false

	currentKv := KeyValueMap{}
	for _, kv := range *kvs {
		// Output always ends with a 'VBoxNetworkName' key.
		if kv.Key == "VBoxNetworkName" && kv.Value != "" {
			ok = true
			if (firstNic == true) && (currentKv["Status"] == "Yes") {
				currentKv["FirstNic"] = "Yes"
				firstNic = true
			}

			kvm[kv.Value] = currentKv
			currentKv = KeyValueMap{}

		} else {
			currentKv[kv.Key] = kv.Value
		}
	}

	return kvm, ok
}

func (kvs *KeyValues) decodeNics() (KeyValuesMap, bool) {

	ok := false
	kvm := make(KeyValuesMap)
	currentName := ""

	currentKv := KeyValueMap{}
	for _, kv := range *kvs {
		// Output always ends with a 'VBoxNetworkName' key.
		if kv.Key == "Name" && kv.Value != "" {
			ok = true
			currentName = kv.Value

			kvm[currentName] = currentKv
			currentKv = KeyValueMap{}

		} else {
			currentKv[kv.Key] = kv.Value
			kvm[currentName] = currentKv
		}
	}

	return kvm, ok
}

func (kvs *KeyValues) decodeShowVmInfo() (KeyValueMap, bool) {

	ok := false
	kvm := make(KeyValueMap)

	for _, kv := range *kvs {
		// Output always start with a 'name' key.
		if kv.Key == "name" && kv.Value != "" {
			ok = true
		}
		if kv.Key != "" && kv.Value != "" {
			kvm[kv.Key] = kv.Value
		}
	}

	return kvm, ok
}
