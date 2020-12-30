package hook

import "github.com/kougazhang/jobx/io"

type Hook struct {
    BeforeReader []func(src io.Src, dst io.Dst) (io.Src, io.Dst, error)
    AfterReader  []func(dst io.Dst) (io.Dst, error)
    Defer        []Defer
}

type Defer struct {
    Func          func(p ...interface{}) error
    Args          []interface{}
    ContinueIfErr bool
}

func ChainDst(dst io.Dst, lst []func(dst io.Dst) (io.Dst, error)) (res io.Dst, err error) {
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

func Chain(src io.Src, dst io.Dst, lst []func(src io.Src, dst io.Dst) (io.Src, io.Dst, error)) (fSrc io.Src, fDst io.Dst, err error) {
    for _, fn := range lst {
        fSrc, fDst, err = fn(src, dst)
        if err != nil {
            return fSrc, fDst, err
        }
        src = fSrc
        dst = fDst
    }
    return fSrc, fDst, err
}

func ChainDefer(lst []Defer) (err error) {
    for _, item := range lst {
        err = item.Func(item.Args...)
        if !item.ContinueIfErr && err != nil {
            break
        }
    }
    return err
}
