package config

import (
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/styles"
)

func RegisterStyle() {
	styles.Register(chroma.MustNewStyle("sqlaid-dark", chroma.StyleEntries{
		chroma.Comment:                  Env.GUI.Editor.Theme.Dark.Comment,
		chroma.CommentHashbang:          Env.GUI.Editor.Theme.Dark.CommentHashbang,
		chroma.CommentMultiline:         Env.GUI.Editor.Theme.Dark.CommentMultiline,
		chroma.CommentPreproc:           Env.GUI.Editor.Theme.Dark.CommentPreproc,
		chroma.CommentSingle:            Env.GUI.Editor.Theme.Dark.CommentSingle,
		chroma.CommentSpecial:           Env.GUI.Editor.Theme.Dark.CommentSpecial,
		chroma.Generic:                  Env.GUI.Editor.Theme.Dark.Generic,
		chroma.GenericDeleted:           Env.GUI.Editor.Theme.Dark.GenericDeleted,
		chroma.GenericEmph:              Env.GUI.Editor.Theme.Dark.GenericEmph,
		chroma.GenericError:             Env.GUI.Editor.Theme.Dark.GenericError,
		chroma.GenericHeading:           Env.GUI.Editor.Theme.Dark.GenericHeading,
		chroma.GenericInserted:          Env.GUI.Editor.Theme.Dark.GenericInserted,
		chroma.GenericOutput:            Env.GUI.Editor.Theme.Dark.GenericOutput,
		chroma.GenericPrompt:            Env.GUI.Editor.Theme.Dark.GenericPrompt,
		chroma.GenericStrong:            Env.GUI.Editor.Theme.Dark.GenericStrong,
		chroma.GenericSubheading:        Env.GUI.Editor.Theme.Dark.GenericSubheading,
		chroma.GenericTraceback:         Env.GUI.Editor.Theme.Dark.GenericTraceback,
		chroma.GenericUnderline:         Env.GUI.Editor.Theme.Dark.GenericUnderline,
		chroma.Error:                    Env.GUI.Editor.Theme.Dark.Error,
		chroma.Keyword:                  Env.GUI.Editor.Theme.Dark.Keyword,
		chroma.KeywordConstant:          Env.GUI.Editor.Theme.Dark.KeywordConstant,
		chroma.KeywordDeclaration:       Env.GUI.Editor.Theme.Dark.KeywordDeclaration,
		chroma.KeywordNamespace:         Env.GUI.Editor.Theme.Dark.KeywordNamespace,
		chroma.KeywordPseudo:            Env.GUI.Editor.Theme.Dark.KeywordPseudo,
		chroma.KeywordReserved:          Env.GUI.Editor.Theme.Dark.KeywordReserved,
		chroma.KeywordType:              Env.GUI.Editor.Theme.Dark.KeywordType,
		chroma.Literal:                  Env.GUI.Editor.Theme.Dark.Literal,
		chroma.LiteralDate:              Env.GUI.Editor.Theme.Dark.LiteralDate,
		chroma.Name:                     Env.GUI.Editor.Theme.Dark.Name,
		chroma.NameAttribute:            Env.GUI.Editor.Theme.Dark.NameAttribute,
		chroma.NameBuiltin:              Env.GUI.Editor.Theme.Dark.NameBuiltin,
		chroma.NameBuiltinPseudo:        Env.GUI.Editor.Theme.Dark.NameBuiltinPseudo,
		chroma.NameClass:                Env.GUI.Editor.Theme.Dark.NameClass,
		chroma.NameConstant:             Env.GUI.Editor.Theme.Dark.NameConstant,
		chroma.NameDecorator:            Env.GUI.Editor.Theme.Dark.NameDecorator,
		chroma.NameEntity:               Env.GUI.Editor.Theme.Dark.NameEntity,
		chroma.NameException:            Env.GUI.Editor.Theme.Dark.NameException,
		chroma.NameFunction:             Env.GUI.Editor.Theme.Dark.NameFunction,
		chroma.NameLabel:                Env.GUI.Editor.Theme.Dark.NameLabel,
		chroma.NameNamespace:            Env.GUI.Editor.Theme.Dark.NameNamespace,
		chroma.NameOther:                Env.GUI.Editor.Theme.Dark.NameOther,
		chroma.NameTag:                  Env.GUI.Editor.Theme.Dark.NameTag,
		chroma.NameVariable:             Env.GUI.Editor.Theme.Dark.NameVariable,
		chroma.NameVariableClass:        Env.GUI.Editor.Theme.Dark.NameVariableClass,
		chroma.NameVariableGlobal:       Env.GUI.Editor.Theme.Dark.NameVariableGlobal,
		chroma.NameVariableInstance:     Env.GUI.Editor.Theme.Dark.NameVariableInstance,
		chroma.LiteralNumber:            Env.GUI.Editor.Theme.Dark.LiteralNumber,
		chroma.LiteralNumberBin:         Env.GUI.Editor.Theme.Dark.LiteralNumberBin,
		chroma.LiteralNumberFloat:       Env.GUI.Editor.Theme.Dark.LiteralNumberFloat,
		chroma.LiteralNumberHex:         Env.GUI.Editor.Theme.Dark.LiteralNumberHex,
		chroma.LiteralNumberInteger:     Env.GUI.Editor.Theme.Dark.LiteralNumberInteger,
		chroma.LiteralNumberIntegerLong: Env.GUI.Editor.Theme.Dark.LiteralNumberIntegerLong,
		chroma.LiteralNumberOct:         Env.GUI.Editor.Theme.Dark.LiteralNumberOct,
		chroma.Operator:                 Env.GUI.Editor.Theme.Dark.Operator,
		chroma.OperatorWord:             Env.GUI.Editor.Theme.Dark.OperatorWord,
		chroma.Other:                    Env.GUI.Editor.Theme.Dark.Other,
		chroma.Punctuation:              Env.GUI.Editor.Theme.Dark.Punctuation,
		chroma.LiteralString:            Env.GUI.Editor.Theme.Dark.LiteralString,
		chroma.LiteralStringBacktick:    Env.GUI.Editor.Theme.Dark.LiteralStringBacktick,
		chroma.LiteralStringChar:        Env.GUI.Editor.Theme.Dark.LiteralStringChar,
		chroma.LiteralStringDoc:         Env.GUI.Editor.Theme.Dark.LiteralStringDoc,
		chroma.LiteralStringDouble:      Env.GUI.Editor.Theme.Dark.LiteralStringDouble,
		chroma.LiteralStringEscape:      Env.GUI.Editor.Theme.Dark.LiteralStringEscape,
		chroma.LiteralStringHeredoc:     Env.GUI.Editor.Theme.Dark.LiteralStringHeredoc,
		chroma.LiteralStringInterpol:    Env.GUI.Editor.Theme.Dark.LiteralStringInterpol,
		chroma.LiteralStringOther:       Env.GUI.Editor.Theme.Dark.LiteralStringOther,
		chroma.LiteralStringRegex:       Env.GUI.Editor.Theme.Dark.LiteralStringRegex,
		chroma.LiteralStringSingle:      Env.GUI.Editor.Theme.Dark.LiteralStringSingle,
		chroma.LiteralStringSymbol:      Env.GUI.Editor.Theme.Dark.LiteralStringSymbol,
		chroma.Text:                     Env.GUI.Editor.Theme.Dark.Text,
		chroma.TextWhitespace:           Env.GUI.Editor.Theme.Dark.TextWhitespace,
		chroma.Background:               Env.GUI.Editor.Theme.Dark.Background,
	}))
	styles.Register(chroma.MustNewStyle("sqlaid-light", chroma.StyleEntries{
		chroma.Comment:                  Env.GUI.Editor.Theme.Light.Comment,
		chroma.CommentHashbang:          Env.GUI.Editor.Theme.Light.CommentHashbang,
		chroma.CommentMultiline:         Env.GUI.Editor.Theme.Light.CommentMultiline,
		chroma.CommentPreproc:           Env.GUI.Editor.Theme.Light.CommentPreproc,
		chroma.CommentSingle:            Env.GUI.Editor.Theme.Light.CommentSingle,
		chroma.CommentSpecial:           Env.GUI.Editor.Theme.Light.CommentSpecial,
		chroma.Generic:                  Env.GUI.Editor.Theme.Light.Generic,
		chroma.GenericDeleted:           Env.GUI.Editor.Theme.Light.GenericDeleted,
		chroma.GenericEmph:              Env.GUI.Editor.Theme.Light.GenericEmph,
		chroma.GenericError:             Env.GUI.Editor.Theme.Light.GenericError,
		chroma.GenericHeading:           Env.GUI.Editor.Theme.Light.GenericHeading,
		chroma.GenericInserted:          Env.GUI.Editor.Theme.Light.GenericInserted,
		chroma.GenericOutput:            Env.GUI.Editor.Theme.Light.GenericOutput,
		chroma.GenericPrompt:            Env.GUI.Editor.Theme.Light.GenericPrompt,
		chroma.GenericStrong:            Env.GUI.Editor.Theme.Light.GenericStrong,
		chroma.GenericSubheading:        Env.GUI.Editor.Theme.Light.GenericSubheading,
		chroma.GenericTraceback:         Env.GUI.Editor.Theme.Light.GenericTraceback,
		chroma.GenericUnderline:         Env.GUI.Editor.Theme.Light.GenericUnderline,
		chroma.Error:                    Env.GUI.Editor.Theme.Light.Error,
		chroma.Keyword:                  Env.GUI.Editor.Theme.Light.Keyword,
		chroma.KeywordConstant:          Env.GUI.Editor.Theme.Light.KeywordConstant,
		chroma.KeywordDeclaration:       Env.GUI.Editor.Theme.Light.KeywordDeclaration,
		chroma.KeywordNamespace:         Env.GUI.Editor.Theme.Light.KeywordNamespace,
		chroma.KeywordPseudo:            Env.GUI.Editor.Theme.Light.KeywordPseudo,
		chroma.KeywordReserved:          Env.GUI.Editor.Theme.Light.KeywordReserved,
		chroma.KeywordType:              Env.GUI.Editor.Theme.Light.KeywordType,
		chroma.Literal:                  Env.GUI.Editor.Theme.Light.Literal,
		chroma.LiteralDate:              Env.GUI.Editor.Theme.Light.LiteralDate,
		chroma.Name:                     Env.GUI.Editor.Theme.Light.Name,
		chroma.NameAttribute:            Env.GUI.Editor.Theme.Light.NameAttribute,
		chroma.NameBuiltin:              Env.GUI.Editor.Theme.Light.NameBuiltin,
		chroma.NameBuiltinPseudo:        Env.GUI.Editor.Theme.Light.NameBuiltinPseudo,
		chroma.NameClass:                Env.GUI.Editor.Theme.Light.NameClass,
		chroma.NameConstant:             Env.GUI.Editor.Theme.Light.NameConstant,
		chroma.NameDecorator:            Env.GUI.Editor.Theme.Light.NameDecorator,
		chroma.NameEntity:               Env.GUI.Editor.Theme.Light.NameEntity,
		chroma.NameException:            Env.GUI.Editor.Theme.Light.NameException,
		chroma.NameFunction:             Env.GUI.Editor.Theme.Light.NameFunction,
		chroma.NameLabel:                Env.GUI.Editor.Theme.Light.NameLabel,
		chroma.NameNamespace:            Env.GUI.Editor.Theme.Light.NameNamespace,
		chroma.NameOther:                Env.GUI.Editor.Theme.Light.NameOther,
		chroma.NameTag:                  Env.GUI.Editor.Theme.Light.NameTag,
		chroma.NameVariable:             Env.GUI.Editor.Theme.Light.NameVariable,
		chroma.NameVariableClass:        Env.GUI.Editor.Theme.Light.NameVariableClass,
		chroma.NameVariableGlobal:       Env.GUI.Editor.Theme.Light.NameVariableGlobal,
		chroma.NameVariableInstance:     Env.GUI.Editor.Theme.Light.NameVariableInstance,
		chroma.LiteralNumber:            Env.GUI.Editor.Theme.Light.LiteralNumber,
		chroma.LiteralNumberBin:         Env.GUI.Editor.Theme.Light.LiteralNumberBin,
		chroma.LiteralNumberFloat:       Env.GUI.Editor.Theme.Light.LiteralNumberFloat,
		chroma.LiteralNumberHex:         Env.GUI.Editor.Theme.Light.LiteralNumberHex,
		chroma.LiteralNumberInteger:     Env.GUI.Editor.Theme.Light.LiteralNumberInteger,
		chroma.LiteralNumberIntegerLong: Env.GUI.Editor.Theme.Light.LiteralNumberIntegerLong,
		chroma.LiteralNumberOct:         Env.GUI.Editor.Theme.Light.LiteralNumberOct,
		chroma.Operator:                 Env.GUI.Editor.Theme.Light.Operator,
		chroma.OperatorWord:             Env.GUI.Editor.Theme.Light.OperatorWord,
		chroma.Other:                    Env.GUI.Editor.Theme.Light.Other,
		chroma.Punctuation:              Env.GUI.Editor.Theme.Light.Punctuation,
		chroma.LiteralString:            Env.GUI.Editor.Theme.Light.LiteralString,
		chroma.LiteralStringBacktick:    Env.GUI.Editor.Theme.Light.LiteralStringBacktick,
		chroma.LiteralStringChar:        Env.GUI.Editor.Theme.Light.LiteralStringChar,
		chroma.LiteralStringDoc:         Env.GUI.Editor.Theme.Light.LiteralStringDoc,
		chroma.LiteralStringDouble:      Env.GUI.Editor.Theme.Light.LiteralStringDouble,
		chroma.LiteralStringEscape:      Env.GUI.Editor.Theme.Light.LiteralStringEscape,
		chroma.LiteralStringHeredoc:     Env.GUI.Editor.Theme.Light.LiteralStringHeredoc,
		chroma.LiteralStringInterpol:    Env.GUI.Editor.Theme.Light.LiteralStringInterpol,
		chroma.LiteralStringOther:       Env.GUI.Editor.Theme.Light.LiteralStringOther,
		chroma.LiteralStringRegex:       Env.GUI.Editor.Theme.Light.LiteralStringRegex,
		chroma.LiteralStringSingle:      Env.GUI.Editor.Theme.Light.LiteralStringSingle,
		chroma.LiteralStringSymbol:      Env.GUI.Editor.Theme.Light.LiteralStringSymbol,
		chroma.Text:                     Env.GUI.Editor.Theme.Light.Text,
		chroma.TextWhitespace:           Env.GUI.Editor.Theme.Light.TextWhitespace,
		chroma.Background:               Env.GUI.Editor.Theme.Light.Background,
	}))
}
