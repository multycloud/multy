package converter

import (
	"github.com/multycloud/multy/resources/output"
	"google.golang.org/protobuf/proto"
)

type ResourceConverters[Arg proto.Message, OutT proto.Message] interface {
	Convert(resourceId string, request Arg, state *output.TfState) (OutT, error)
}
