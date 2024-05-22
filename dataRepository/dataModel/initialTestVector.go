package dataModel

type InitialTestVevtor struct{}

func (InitialTestVevtor) Patients() []Patient {
	return []Patient{
		{Name: "Patient 1", Occupation: "Sw Developer",
			City: "Krakow Krowodrza", TelephoneNumber: "+486111111111", BirthYear: 1982},

		{Name: "Patient 2", Occupation: "Student",
			City: "Krakow Pradnik", TelephoneNumber: "+48222222222", BirthYear: 2007},

		{Name: "Patient 3", Occupation: "Therapist",
			City: "Krakow Bronowice", TelephoneNumber: "+48760300300", BirthYear: 1979},
	}
}

func (InitialTestVevtor) Notes() []Note {
	return []Note{
		{Name: "Note 1", PatientRowId: 1,
			SessionDate: 1, NoteDate: 1, FileName: "", IsCrypted: false},
		{Name: "Note 2", PatientRowId: 1,
			SessionDate: 1, NoteDate: 1, FileName: "", IsCrypted: false},
		{Name: "Note 3", PatientRowId: 1,
			SessionDate: 1, NoteDate: 1, FileName: "", IsCrypted: false},
	}
}

func (InitialTestVevtor) Users() []User {
	return []User{
		{Name: "User", LastName: "One", Email: "fake1@gmail.com",
			TelephoneNumber: "+486111111111", Salt: "447f44a4", PubKey: "",
			Password: "b1ecb61a0c76f7bbb253d80c6d610818ceb2cdfbc16b923036a93a61229d426e"}, // aslk12
		{Name: "User", LastName: "Two", Email: "fake2@gmail.com",
			TelephoneNumber: "+486222222222", Salt: "447f44a4", PubKey: "",
			Password: "b1ecb61a0c76f7bbb253d80c6d610818ceb2cdfbc16b923036a93a61229d426e"}, // aslk12

	}

}

func (InitialTestVevtor) Manifests() []PatientManifest {

	return []PatientManifest{
		{PatientId: 1, UserId: 1, CrudMask: 0x7, EncryptedAes: "test"},
	}

}
