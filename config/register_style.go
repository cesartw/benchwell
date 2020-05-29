package config

import (
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/styles"
)

func RegisterStyle() {
	styles.Register(chroma.MustNewStyle("sqlaid-dark", chroma.StyleEntries{
		chroma.Comment:                  Env.GUI.Editor.Theme.Comment,
		chroma.CommentHashbang:          Env.GUI.Editor.Theme.CommentHashbang,
		chroma.CommentMultiline:         Env.GUI.Editor.Theme.CommentMultiline,
		chroma.CommentPreproc:           Env.GUI.Editor.Theme.CommentPreproc,
		chroma.CommentSingle:            Env.GUI.Editor.Theme.CommentSingle,
		chroma.CommentSpecial:           Env.GUI.Editor.Theme.CommentSpecial,
		chroma.Generic:                  Env.GUI.Editor.Theme.Generic,
		chroma.GenericDeleted:           Env.GUI.Editor.Theme.GenericDeleted,
		chroma.GenericEmph:              Env.GUI.Editor.Theme.GenericEmph,
		chroma.GenericError:             Env.GUI.Editor.Theme.GenericError,
		chroma.GenericHeading:           Env.GUI.Editor.Theme.GenericHeading,
		chroma.GenericInserted:          Env.GUI.Editor.Theme.GenericInserted,
		chroma.GenericOutput:            Env.GUI.Editor.Theme.GenericOutput,
		chroma.GenericPrompt:            Env.GUI.Editor.Theme.GenericPrompt,
		chroma.GenericStrong:            Env.GUI.Editor.Theme.GenericStrong,
		chroma.GenericSubheading:        Env.GUI.Editor.Theme.GenericSubheading,
		chroma.GenericTraceback:         Env.GUI.Editor.Theme.GenericTraceback,
		chroma.GenericUnderline:         Env.GUI.Editor.Theme.GenericUnderline,
		chroma.Error:                    Env.GUI.Editor.Theme.Error,
		chroma.Keyword:                  Env.GUI.Editor.Theme.Keyword,
		chroma.KeywordConstant:          Env.GUI.Editor.Theme.KeywordConstant,
		chroma.KeywordDeclaration:       Env.GUI.Editor.Theme.KeywordDeclaration,
		chroma.KeywordNamespace:         Env.GUI.Editor.Theme.KeywordNamespace,
		chroma.KeywordPseudo:            Env.GUI.Editor.Theme.KeywordPseudo,
		chroma.KeywordReserved:          Env.GUI.Editor.Theme.KeywordReserved,
		chroma.KeywordType:              Env.GUI.Editor.Theme.KeywordType,
		chroma.Literal:                  Env.GUI.Editor.Theme.Literal,
		chroma.LiteralDate:              Env.GUI.Editor.Theme.LiteralDate,
		chroma.Name:                     Env.GUI.Editor.Theme.Name,
		chroma.NameAttribute:            Env.GUI.Editor.Theme.NameAttribute,
		chroma.NameBuiltin:              Env.GUI.Editor.Theme.NameBuiltin,
		chroma.NameBuiltinPseudo:        Env.GUI.Editor.Theme.NameBuiltinPseudo,
		chroma.NameClass:                Env.GUI.Editor.Theme.NameClass,
		chroma.NameConstant:             Env.GUI.Editor.Theme.NameConstant,
		chroma.NameDecorator:            Env.GUI.Editor.Theme.NameDecorator,
		chroma.NameEntity:               Env.GUI.Editor.Theme.NameEntity,
		chroma.NameException:            Env.GUI.Editor.Theme.NameException,
		chroma.NameFunction:             Env.GUI.Editor.Theme.NameFunction,
		chroma.NameLabel:                Env.GUI.Editor.Theme.NameLabel,
		chroma.NameNamespace:            Env.GUI.Editor.Theme.NameNamespace,
		chroma.NameOther:                Env.GUI.Editor.Theme.NameOther,
		chroma.NameTag:                  Env.GUI.Editor.Theme.NameTag,
		chroma.NameVariable:             Env.GUI.Editor.Theme.NameVariable,
		chroma.NameVariableClass:        Env.GUI.Editor.Theme.NameVariableClass,
		chroma.NameVariableGlobal:       Env.GUI.Editor.Theme.NameVariableGlobal,
		chroma.NameVariableInstance:     Env.GUI.Editor.Theme.NameVariableInstance,
		chroma.LiteralNumber:            Env.GUI.Editor.Theme.LiteralNumber,
		chroma.LiteralNumberBin:         Env.GUI.Editor.Theme.LiteralNumberBin,
		chroma.LiteralNumberFloat:       Env.GUI.Editor.Theme.LiteralNumberFloat,
		chroma.LiteralNumberHex:         Env.GUI.Editor.Theme.LiteralNumberHex,
		chroma.LiteralNumberInteger:     Env.GUI.Editor.Theme.LiteralNumberInteger,
		chroma.LiteralNumberIntegerLong: Env.GUI.Editor.Theme.LiteralNumberIntegerLong,
		chroma.LiteralNumberOct:         Env.GUI.Editor.Theme.LiteralNumberOct,
		chroma.Operator:                 Env.GUI.Editor.Theme.Operator,
		chroma.OperatorWord:             Env.GUI.Editor.Theme.OperatorWord,
		chroma.Other:                    Env.GUI.Editor.Theme.Other,
		chroma.Punctuation:              Env.GUI.Editor.Theme.Punctuation,
		chroma.LiteralString:            Env.GUI.Editor.Theme.LiteralString,
		chroma.LiteralStringBacktick:    Env.GUI.Editor.Theme.LiteralStringBacktick,
		chroma.LiteralStringChar:        Env.GUI.Editor.Theme.LiteralStringChar,
		chroma.LiteralStringDoc:         Env.GUI.Editor.Theme.LiteralStringDoc,
		chroma.LiteralStringDouble:      Env.GUI.Editor.Theme.LiteralStringDouble,
		chroma.LiteralStringEscape:      Env.GUI.Editor.Theme.LiteralStringEscape,
		chroma.LiteralStringHeredoc:     Env.GUI.Editor.Theme.LiteralStringHeredoc,
		chroma.LiteralStringInterpol:    Env.GUI.Editor.Theme.LiteralStringInterpol,
		chroma.LiteralStringOther:       Env.GUI.Editor.Theme.LiteralStringOther,
		chroma.LiteralStringRegex:       Env.GUI.Editor.Theme.LiteralStringRegex,
		chroma.LiteralStringSingle:      Env.GUI.Editor.Theme.LiteralStringSingle,
		chroma.LiteralStringSymbol:      Env.GUI.Editor.Theme.LiteralStringSymbol,
		chroma.Text:                     Env.GUI.Editor.Theme.Text,
		chroma.TextWhitespace:           Env.GUI.Editor.Theme.TextWhitespace,
		chroma.Background:               Env.GUI.Editor.Theme.Background,
	}))
}
