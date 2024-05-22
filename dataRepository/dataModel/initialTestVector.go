package dataModel

type InitialTestVevtor struct {
}

func (InitialTestVevtor) Patients() []Patient {
	pateintsTestVector :=
		[]Patient{
			{Name: "Patient 1", Occupation: "Sw Developer",
				City: "Krakow Krowodrza", TelephoneNumber: "+486111111111", BirthYear: 1982},

			{Name: "Patient 2", Occupation: "Student",
				City: "Krakow Pradnik", TelephoneNumber: "+48222222222", BirthYear: 2007},

			{Name: "Patient 3", Occupation: "Therapist",
				City: "Krakow Bronowice", TelephoneNumber: "+48760300300", BirthYear: 1979},
		}
	return pateintsTestVector
}

func (InitialTestVevtor) Notes() []Note {
	notesTestVector := []Note{
		{Name: "Note 1", PatientRowId: 1,
			SessionDate: 1, NoteDate: 1, FileName: "test1.txt", IsCrypted: false},
		{Name: "Note 2", PatientRowId: 1,
			SessionDate: 1, NoteDate: 1, FileName: "test2.txt", IsCrypted: false},
		{Name: "Note 3", PatientRowId: 1,
			SessionDate: 1, NoteDate: 1, FileName: "test3.txt", IsCrypted: false},
	}
	return notesTestVector
}

func (InitialTestVevtor) Users() []User {
	return []User{
		{Name: "User", LastName: "1", Email: "briliantFakeUser@gmail.con",
			TelephoneNumber: "+486111111111", PasswordSalt: "DeadBeef", PasswordSha: "", PubKey: ""},
	}
}

func (InitialTestVevtor) Manifests() []PatientManifest {

	return []PatientManifest{
		{PatientId: 1, UserId: 1, CrudMask: 0x7, EncryptedAes: "test"},
	}

}
