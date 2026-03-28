package flags

import "flag"

type Flags struct {
    Path  string
    Debug bool
}

func Parse() *Flags {
    config := flag.String("config", "./config/config.yaml", "path to config")
    debug := flag.Bool("debug", false, "enable debug mode")
    flag.Parse()
    
    return &Flags{
        Path:  *config,
        Debug: *debug,
    }
}
