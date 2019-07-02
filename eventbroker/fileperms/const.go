package fileperms

// UNIX file permission constants.
const (
	Read       = 04
	Write      = 02
	Execute    = 01
	UserShift  = 6
	GroupShift = 3
	OtherShift = 0

	UserRead             = Read << UserShift
	UserWrite            = Write << UserShift
	UserExecute          = Execute << UserShift
	UserReadWrite        = UserRead | UserWrite
	UserReadWriteExecute = UserReadWrite | UserExecute

	GroupRead             = Read << GroupShift
	GroupWrite            = Write << GroupShift
	GroupExecute          = Execute << GroupShift
	GroupReadWrite        = GroupRead | GroupWrite
	GroupReadWriteExecute = GroupReadWrite | GroupExecute

	OtherRead             = Read << OtherShift
	OtherWrite            = Write << OtherShift
	OtherExecute          = Execute << OtherShift
	OtherReadWrite        = OtherRead | OtherWrite
	OtherReadWriteExecute = OtherReadWrite | OtherExecute

	AllRead             = UserRead | GroupRead | OtherRead
	AllWrite            = UserWrite | GroupWrite | OtherWrite
	AllExecute          = UserExecute | GroupExecute | OtherExecute
	AllReadWrite        = AllRead | AllWrite
	AllReadWriteExecute = AllReadWrite | GroupExecute
)
