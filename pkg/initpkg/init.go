package initpkg

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/itchyny/gojq"
)

func GetGoModPath() (string, error) {
	goModInfo, err := exec.Command("go", "mod", "edit", "-json").Output()

	fmt.Println(err,"===========err=========")

	if err != nil {
		return "", err
	}

	query, err := gojq.Parse(".Module.Path | ..")

	fmt.Println(err,"===========err=========")

	if err != nil {
		return "", err
	}

	goModDetails := make(map[string]interface{})

	errJson := json.Unmarshal(goModInfo, &goModDetails)
	if errJson != nil {
		log.Fatal(errJson)
	}

	iter := query.Run(goModDetails)

	var goModPath string
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			log.Fatal(err)
		}
		goModPath = v.(string)
	}

	fmt.Println(goModPath,"===========gomodpath=========")

	return goModPath, nil
}
