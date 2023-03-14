package telegraph

import (
	"strings"

	telegraphPkg "github.com/anonyindian/telegraph-go"
	"go.uber.org/multierr"
)

func UploadFiles(pPaths []string) ([]string, error) {
	var uploaded = []string{}
	var err error
	for _, pPath := range pPaths {
		url, err1 := telegraphPkg.UploadFile(pPath)
		if err != nil {
			err = multierr.Append(err, err1)
			continue
		}
		uploaded = append(uploaded, strings.Join([]string{"telegra.ph", url}, ""))
	}
	return uploaded, err
}
