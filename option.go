package gotts

/**
 * @Author zyq
 * @Date 2023/2/21 3:38 PM
 * @Description
 **/

type Option struct {
    w            Writer
    token, voice string
    module       ConnModule
}
type Options func(option *Option)

func GetOption(ops ...Options) *Option {
    option := new(Option)
    for _, v := range ops {
        v(option)
    }
    if option.module == "" {
        option.module = ConnModuleBing
    }
    return option
}

func WithWriter(w Writer) Options {
    return func(option *Option) {
        option.w = w
    }
}

func WithVoice(voice string) Options {
    return func(option *Option) {
        option.voice = voice
    }
}

func WithToken(token string) Options {
    return func(option *Option) {
        option.token = token
    }
}

func WithModule(module ConnModule) Options {
    return func(option *Option) {
        option.module = module
    }
}
