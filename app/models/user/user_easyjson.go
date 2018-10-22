// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package user

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson9e1087fdDecodeTechDbServerAppModelsUser(in *jlexer.Lexer, out *UserUpdate) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "about":
			if m, ok := out.About.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.About.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.About = in.Interface()
			}
		case "email":
			if m, ok := out.Email.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Email.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Email = in.Interface()
			}
		case "fullname":
			if m, ok := out.Fullname.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Fullname.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Fullname = in.Interface()
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9e1087fdEncodeTechDbServerAppModelsUser(out *jwriter.Writer, in UserUpdate) {
	out.RawByte('{')
	first := true
	_ = first
	if in.About != nil {
		const prefix string = ",\"about\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if m, ok := in.About.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.About.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.About))
		}
	}
	{
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if m, ok := in.Email.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Email.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Email))
		}
	}
	{
		const prefix string = ",\"fullname\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if m, ok := in.Fullname.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Fullname.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Fullname))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserUpdate) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeTechDbServerAppModelsUser(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserUpdate) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeTechDbServerAppModelsUser(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserUpdate) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeTechDbServerAppModelsUser(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserUpdate) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeTechDbServerAppModelsUser(l, v)
}
func easyjson9e1087fdDecodeTechDbServerAppModelsUser1(in *jlexer.Lexer, out *UserPointList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(UserPointList, 0, 8)
			} else {
				*out = UserPointList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 *User
			if in.IsNull() {
				in.Skip()
				v1 = nil
			} else {
				if v1 == nil {
					v1 = new(User)
				}
				(*v1).UnmarshalEasyJSON(in)
			}
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9e1087fdEncodeTechDbServerAppModelsUser1(out *jwriter.Writer, in UserPointList) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			if v3 == nil {
				out.RawString("null")
			} else {
				(*v3).MarshalEasyJSON(out)
			}
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v UserPointList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeTechDbServerAppModelsUser1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserPointList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeTechDbServerAppModelsUser1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserPointList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeTechDbServerAppModelsUser1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserPointList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeTechDbServerAppModelsUser1(l, v)
}
func easyjson9e1087fdDecodeTechDbServerAppModelsUser2(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "about":
			out.About = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "fullname":
			out.Fullname = string(in.String())
		case "nickname":
			out.Nickname = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9e1087fdEncodeTechDbServerAppModelsUser2(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	if in.About != "" {
		const prefix string = ",\"about\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.About))
	}
	{
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"fullname\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Fullname))
	}
	if in.Nickname != "" {
		const prefix string = ",\"nickname\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Nickname))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeTechDbServerAppModelsUser2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeTechDbServerAppModelsUser2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeTechDbServerAppModelsUser2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeTechDbServerAppModelsUser2(l, v)
}
