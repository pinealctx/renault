package paths

import "os"

func Exist(path string) (bool, error) {
	var _, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
