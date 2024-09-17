package models

import (
	"ms-sv-jira/helper"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// BRIBrainClaims for JWT token
type BRIBrainClaims struct {
	User
	ExpiredToken
	jwt.StandardClaims
}

type ExpiredToken struct {
	ExpiredToken int64 `json:"expired_token,omitempty"`
}

type User struct {
	Pernr        string  `json:"pernr" gorm:"column:PERNR"`
	Sname        string  `json:"sname" gorm:"column:SNAME"`
	Role         string  `json:"role" gorm:"column:ROLE"`
	Description1 string  `json:"description_1" gorm:"column:DESCRIPTION_1"`
	Description2 string  `json:"description_2" gorm:"column:DESCRIPTION_2"`
	OrgehTx      string  `json:"orgeh_tx" gorm:"column:ORGEH_TX"`
	StellTx      string  `json:"stell_tx" gorm:"column:STELL_TX"`
	Rgdesc       string  `json:"rgdesc" gorm:"column:RGDESC"`
	Mbdesc       string  `json:"mbdesc" gorm:"column:MBDESC"`
	Brdesc       string  `json:"brdesc" gorm:"column:BRDESC"`
	Region       string  `json:"region" gorm:"column:REGION"`
	Mainbranch   float64 `json:"mainbr" gorm:"column:MAINBR"`
	Branch       float64 `json:"branch" gorm:"column:BRANCH"`
	JenisKelamin string  `json:"jenis_kelamin" gorm:"column:JENIS_KELAMIN"`
	AccessLevel  string  `json:"access_level" gorm:"column:ACCESS_LEVEL"`
	PnKecamatan  string  `json:"pn_kecamatan" gorm:"column:PN_KECAMATAN"`
}
type BribrainAuditTrailRequest struct {
	User    string `json:"user"`
	Role    string `json:"role"`
	Service string `json:"service"`
	Action  string `json:"action"`
	Remark  string `json:"remark"`
	Payload string `json:"payload"`
}

type PayloadBribrainAuditTrail struct {
	Url          string                           `json:"url"`
	Headers      HeadersPayloadBribrainAuditTrail `json:"headers"`
	BodyRequest  string                           `json:"body_request"`
	BodyResponse string                           `json:"body_response"`
	StatusCode   int                              `json:"status_code"`
}
type HeadersPayloadBribrainAuditTrail struct {
	ContentType   string `json:"content_type"`
	Authorization string `json:"authorization"`
}
type ValidateAccess struct {
	Roles       []string `json:"roles"`
	AccessLevel []string `json:"access_level"`
}

func (BribrainAuditTrailRequest) MappingBribrainAuditTrail(action string, remark string, user, role string, payload string) BribrainAuditTrailRequest {
	res := BribrainAuditTrailRequest{
		User:    user,
		Role:    role,
		Service: os.Getenv("APP_NAME"),
		Action:  action,
		Remark:  remark,
		Payload: payload,
	}
	return res
}
func (PayloadBribrainAuditTrail) MappingPayloadBribrainAuditTrail(url string, headers HeadersPayloadBribrainAuditTrail, bodyrequest, bodyresponse string, statuscode int) PayloadBribrainAuditTrail {
	res := PayloadBribrainAuditTrail{
		Url:          url,
		Headers:      headers,
		BodyRequest:  bodyrequest,
		BodyResponse: bodyresponse,
		StatusCode:   statuscode,
	}
	return res
}
func (HeadersPayloadBribrainAuditTrail) MappingHeadersPayloadBribrainAuditTrail(contentype string, authorization string) HeadersPayloadBribrainAuditTrail {
	res := HeadersPayloadBribrainAuditTrail{
		ContentType:   contentype,
		Authorization: authorization,
	}
	return res
}
func (claims *BRIBrainClaims) MappingToGeneralRequestV3() *GeneralRequestV3 {
	req := &GeneralRequestV3{}
	if claims.Role == "MANTRI" {
		req.Region = claims.Region
		req.Mainbranch = helper.FloatToString(claims.Mainbranch)
		req.Branch = claims.Branch
		req.Role = claims.Role
		req.Mantri = claims.Pernr
	} else if claims.Role == "KP" || claims.Role == "KW" || claims.Role == "KC" {
		req.Role = claims.Role
	} else {
		req.Role = "INVALID"
	}
	return req
}

type GeneralRequestV3 struct {
	// Kategori              string    `json:"kategori" validate:"required,required"`
	MID                   string    `json:"mid"`
	IDNumber              string    `json:"id_number"`
	KategoriRekomendasi   string    `json:"kategori" validate:"required,required"`
	Periode               string    `json:"periode"`
	Region                string    `json:"region" validate:"required,required"`
	Mainbranch            string    `json:"mainbranch" validate:"required,required"`
	Branch                float64   `json:"branch" validate:"required,required"`
	Role                  string    `json:"role" validate:"required,required"`
	Pn                    string    `json:"pn" validate:"required"`
	Mantri                string    `json:"mantri" validate:"required,required"`
	KategoriSkor          string    `json:"kategori_skor" validate:"required,required"`
	Limit                 int       `json:"limit"`
	Page                  int       `json:"page"`
	Search                string    `json:"search"`
	Acctno                int       `json:"acctno"`
	DeskripsiTindakLanjut string    `json:"deskripsi_tindak_lanjut"`
	HasilRekomendasi      string    `json:"hasil_rekomendasi"`
	Sort                  string    `json:"sort"`
	Asc                   string    `json:"asc"`
	FlagTab               int       `json:"flag_tab"`
	FlagTindakLanjut      string    `json:"flag_tindak_lanjut"`
	KodeTLHubungi         float64   `json:"kode_tl_hubungi"`
	KodeTLKunjungi        float64   `json:"kode_tl_kunjungi"`
	KodeKetertarikan      float64   `json:"kode_ketertarikan"`
	TanggalHubungi        string    `json:"tanggal_hubungi"`
	TanggalKunjungi       string    `json:"tanggal_kunjungi"`
	KodeRangeGaji         float64   `json:"kode_range_gaji"`
	FilterTL              []float64 `json:"filter_tl"`
	KodeTL                float64   `json:"kode_tl"`
	NomorSurat            string    `json:"nomor_surat"`
	AuthHeader            string    `json:"auth_header"`
}
