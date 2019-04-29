package api

import "gearbox/apiworks"

const (
	Port        = "9999"
	ServiceName = "Gearbox API"
	Version     = "0.1"
	DocsUrl     = "https://docs.gearbox.works/api"
)

const (
	BaseUrlPattern               Link     = "http://127.0.0.1:%s"
	GearboxApiSchema             Link     = "https://docs.gearbox.works/api/schema/1.0/"
	GearboxApiIdentifier         Metaname = "GearboxAPI"
	MetaGearboxBaseurl                    = GearboxApiIdentifier + ".baseurl"
	MetaGearboxApiSchema                  = GearboxApiIdentifier + ".schema"
	SchemaGearboxApiRelationType          = RelType("schema." + GearboxApiIdentifier)
)

const (
	RelTypePattern = apiworks.RelTypePattern
	SelfRelType    = apiworks.SelfRelType
	CurrentRelType = apiworks.CurrentRelType
	RootRelType    = apiworks.RootRelType
	ItemRelType    = apiworks.ItemRelType
	ListRelType    = apiworks.ListRelType
	AddItemRelType = apiworks.AddItemRelType

	DefaultLanguage = apiworks.DefaultLanguage
	DcTermsSchema   = apiworks.DcTermsSchema
	DcSchema        = apiworks.DcSchema

	SchemaDcRelType       = apiworks.SchemaDcRelType
	SchemaDcTermsRelType  = apiworks.SchemaDcTermsRelType
	MetaDcCreator         = apiworks.MetaDcCreator
	MetaDcTermsIdentifier = apiworks.MetaDcTermsIdentifier
	MetaDcLanguage        = apiworks.MetaDcLanguage

	CharsetUTF8  = apiworks.CharsetUTF8
	NoFilterPath = apiworks.NoFilterPath
	Basepath     = apiworks.Basepath
)
