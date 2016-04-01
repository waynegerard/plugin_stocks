package stocks

import (
        "log"

        "github.com/itsabot/abot/shared/datatypes"
        "github.com/itsabot/abot/shared/nlp"
        "github.com/itsabot/abot/shared/plugin"
)

var p *dt.Plugin

func init() {
    trigger := &nlp.StructuredInput{
        Commands: []string{"find", "show"},
        Objects: []string{"stocks", "portfolio", "etf"},
    }

    p.Vocab = dt.NewVocab(
        dt.VocabHandler{
            Fn: kwFindStocks,
            Trigger: &nlp.StructuredInput{
                Commands: []string{"find"},
                Objects: []string{"stocks", "etf"},
            },
        },
    )

    fns := &dt.PluginFns{Run: Run, FollowUp: FollowUp}

    var err error
    pluginPath := "github.com/wgerard/plugin_stocks"
    p, err = plugin.New(pluginPath, trigger, fns)
    if err != nil {
        log.Fatalln("building", err)
    }
}

func kwFindStocks(in *dt.Msg) (resp string) {
    var s *stock.Stock
    for _, obj := range in.StructuredInput.Objects {
        // Ticker symbols can be at most 5 letters long
        if len(obj) > 5 {
            continue
        }
        s = stock.Get(obj)
    }
    // In the case of no stock found, return an error message.
    // If you return an empty string, Abot will respond to the
    // user with confusion, like "I'm not sure what you mean."
    if s == nil {
        return "I couldn't find any stocks like that."
    }
    return s.Name + " is trading at " + s.Price
}



func Run(in *dt.Msg) (string, error) {
        return FollowUp(in)
}

func FollowUp(in *dt.Msg) (string, error) {
        return p.Vocab.HandleKeywords(in), nil
}

func er(err error) string {
        p.Log.Debug(err)
        return "Something went wrong, but I'll try to get that fixed right away."
}
