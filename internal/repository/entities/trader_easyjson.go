// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package entities

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

func easyjson1339d61cDecodeGithubComBurnbSignallerInternalRepositoryEntities(in *jlexer.Lexer, out *Trader) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "encryptedUid":
			out.Uid = string(in.String())
		case "pnlValue":
			out.Pnl = float64(in.Float64())
		case "weeklyPnl":
			out.PnlWeekly = float64(in.Float64())
		case "monthlyPnl":
			out.PnlMonthly = float64(in.Float64())
		case "yearlyPnl":
			out.PnlYearly = float64(in.Float64())
		case "roiValue":
			out.Roi = float64(in.Float64())
		case "weeklyRoi":
			out.RoiWeekly = float64(in.Float64())
		case "monthlyRoi":
			out.RoiMonthly = float64(in.Float64())
		case "yearlyRoi":
			out.RoiYearly = float64(in.Float64())
		case "positionShared":
			out.PositionShared = bool(in.Bool())
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
func easyjson1339d61cEncodeGithubComBurnbSignallerInternalRepositoryEntities(out *jwriter.Writer, in Trader) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"encryptedUid\":"
		out.RawString(prefix[1:])
		out.String(string(in.Uid))
	}
	{
		const prefix string = ",\"pnlValue\":"
		out.RawString(prefix)
		out.Float64(float64(in.Pnl))
	}
	{
		const prefix string = ",\"weeklyPnl\":"
		out.RawString(prefix)
		out.Float64(float64(in.PnlWeekly))
	}
	{
		const prefix string = ",\"monthlyPnl\":"
		out.RawString(prefix)
		out.Float64(float64(in.PnlMonthly))
	}
	{
		const prefix string = ",\"yearlyPnl\":"
		out.RawString(prefix)
		out.Float64(float64(in.PnlYearly))
	}
	{
		const prefix string = ",\"roiValue\":"
		out.RawString(prefix)
		out.Float64(float64(in.Roi))
	}
	{
		const prefix string = ",\"weeklyRoi\":"
		out.RawString(prefix)
		out.Float64(float64(in.RoiWeekly))
	}
	{
		const prefix string = ",\"monthlyRoi\":"
		out.RawString(prefix)
		out.Float64(float64(in.RoiMonthly))
	}
	{
		const prefix string = ",\"yearlyRoi\":"
		out.RawString(prefix)
		out.Float64(float64(in.RoiYearly))
	}
	{
		const prefix string = ",\"positionShared\":"
		out.RawString(prefix)
		out.Bool(bool(in.PositionShared))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Trader) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1339d61cEncodeGithubComBurnbSignallerInternalRepositoryEntities(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Trader) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson1339d61cEncodeGithubComBurnbSignallerInternalRepositoryEntities(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Trader) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson1339d61cDecodeGithubComBurnbSignallerInternalRepositoryEntities(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Trader) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson1339d61cDecodeGithubComBurnbSignallerInternalRepositoryEntities(l, v)
}
