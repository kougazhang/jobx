package io

type Src interface{}

type Dst interface{}

type IO interface {
	Copy(src Src, dst Dst) error
}

func ChainDst(dst Dst, lst []func(dst Dst) (Dst, error)) (res Dst, err error) {
	if len(lst) == 0 {
		return dst, err
	}

	for _, fn := range lst {
		res, err = fn(dst)
		if err != nil {
			return res, err
		}
		dst = res
	}
	return res, err
}

func ChainSrc(src Src, lst []func(src Src) (Dst, error)) (res Src, err error) {
	for _, fn := range lst {
		res, err = fn(src)
		if err != nil {
			return res, err
		}
		src = res
	}
	return res, err
}
