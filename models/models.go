package models

import (
	"sync"
	"time"

	"gorm.io/gorm/schema"
)

type HcdListAndDetailEmployee struct {
	Pn                           string
	Nama                         string
	Level                        string
	NomorTelepon                 string
	JenisKelamin                 string
	CorporateTitle               string
	JabatanTerakhir              string
	JabatanTerakhirAssessment    string
	KodeArea                     string
	Area                         string
	KodeUker                     string
	UnitKerja                    string
	UnitKerjaAssessment          string
	DigitalLeadership            int
	GlobalBusinessSavvy          int
	CustomerFocus                int
	BuildingStrategicPartnership int
	StrategicOrientation         int
	DrivingExecution             int
	DrivingInnovation            int
	DevelopingOrganizations      int
	LeadingChange                int
	ManagingDiversity            int
	PqOperations                 int
	PqFinance                    int
	PqPeople                     int
	PqTechnology                 int
	PqCommercial                 int
	PenyelenggaraAssessment      string
	RekomendasiBri               string
	RekomendasiBumn              string
	NilaiKompetensiBumn          float64
	TanggalAssessment            time.Time
}

func (*HcdListAndDetailEmployee) Parse() *HcdListAndDetailEmployee {
	s, _ := schema.Parse(&HcdListAndDetailEmployee{}, &sync.Map{}, schema.NamingStrategy{
		TablePrefix:   "bribrain_",
		SingularTable: true,
	})
	var column []string
	for _, field := range s.Fields {
		col := field.DBName
		column = append(column, col)
	}
	return &HcdListAndDetailEmployee{
		Pn:              s.Table + "." + column[0],
		Nama:            s.Table + "." + column[1],
		Level:           s.Table + "." + column[2],
		CorporateTitle:  s.Table + "." + column[5],
		JabatanTerakhir: s.Table + "." + column[6],
		KodeArea:        s.Table + "." + column[8],
		Area:            s.Table + "." + column[9],
		UnitKerja:       s.Table + "." + column[11],
	}
}
