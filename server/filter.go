package server

import (
	"fmt"

	"github.com/teacat/jsonfilter"
)

func ApplyResourceListFilter(data []byte) ([]byte, error) {
	return jsonfilter.Filter(data, "items(apiVersion,kind,metadata(annotations,creationTimestamp,finalizers,generation,name,namespace,resourceVersion,uid),spec,status)")
}

func ApplyResourceFilter(format string, data []byte) ([]byte, error) {
	return jsonfilter.Filter(data, fmt.Sprintf(format, "apiVersion,kind,metadata(annotations,creationTimestamp,finalizers,generation,name,namespace,resourceVersion,uid),spec,status"))
}
