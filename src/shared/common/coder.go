package common

import "github.com/ServiceWeaver/weaver/runtime/codegen"

func Encode_slice_int64(enc *codegen.Encoder, slice []int64) {
	if slice == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(slice))
	for i := 0; i < len(slice); i++ {
		enc.Int64(slice[i])
	}
}

func Encode_slice_string(enc *codegen.Encoder, slice []string) {
	if slice == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(slice))
	for i := 0; i < len(slice); i++ {
		enc.String(slice[i])
	}
}

func Decode_slice_int64(dec *codegen.Decoder) []int64 {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]int64, n)
	for i := 0; i < n; i++ {
		res[i] = dec.Int64()
	}
	return res
}

func Decode_slice_string(dec *codegen.Decoder) []string {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		res[i] = dec.String()
	}
	return res
}
