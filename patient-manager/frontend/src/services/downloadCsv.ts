import { getAllPatients, getCheckupsForRecord, getIllnessesForRecord, getPrescriptionsForIllness } from './patientService';

export async function downloadAllPatientsData(): Promise<void> {
    const patients = await getAllPatients();

    if (!patients || patients.length === 0) {
        throw new Error('No patients found to download.');
    }

    const allData = [];

    const headers = [
        "Record Type", "Patient First Name", "Patient Last Name", "Patient OIB",
        "Checkup UUID", "Checkup Date", "Checkup Type", "Checkup Illness ID",
        "Illness UUID", "Illness Name", "Illness Start Date", "Illness End Date",
        "Prescription UUID", "Prescription Issued At", "Prescription Illness Name", "Medications"
    ];

    allData.push(headers);

    for (const patient of patients) {
        const patientInfo = {
            "Record Type": "Patient",
            "Patient First Name": patient.firstName,
            "Patient Last Name": patient.lastName,
            "Patient OIB": patient.oib,
            "Checkup UUID": "", "Checkup Date": "", "Checkup Type": "", "Checkup Illness ID": "",
            "Illness UUID": "", "Illness Name": "", "Illness Start Date": "", "Illness End Date": "",
            "Prescription UUID": "", "Prescription Issued At": "", "Prescription Illness Name": "", "Medications": ""
        };
        
        allData.push(Object.values(patientInfo));

        const checkups = await getCheckupsForRecord(patient.medicalRecordUuid);
        const illnesses = await getIllnessesForRecord(patient.medicalRecordUuid);
        
        let allPrescriptions = [];
        for (const illness of illnesses) {
            const prescriptions = await getPrescriptionsForIllness(illness.id);
            if (prescriptions) {
                allPrescriptions.push(...prescriptions.map(p => ({
                    ...p,
                    illnessName: illness.name
                })));
            }
        }

        for (const checkup of checkups) {
            const checkupRow = {
                "Record Type": "Checkup",
                "Patient First Name": patient.firstName,
                "Patient Last Name": patient.lastName,
                "Patient OIB": patient.oib,
                "Checkup UUID": checkup.uuid,
                "Checkup Date": new Date(checkup.checkupDate).toLocaleDateString(),
                "Checkup Type": checkup.type,
                "Checkup Illness ID": checkup.illnessId || "N/A",
                "Illness UUID": "", "Illness Name": "", "Illness Start Date": "", "Illness End Date": "",
                "Prescription UUID": "", "Prescription Issued At": "", "Prescription Illness Name": "", "Medications": ""
            };
            allData.push(Object.values(checkupRow));
        }

        for (const illness of illnesses) {
            const illnessRow = {
                "Record Type": "Illness",
                "Patient First Name": patient.firstName,
                "Patient Last Name": patient.lastName,
                "Patient OIB": patient.oib,
                "Checkup UUID": "", "Checkup Date": "", "Checkup Type": "", "Checkup Illness ID": "",
                "Illness UUID": illness.uuid,
                "Illness Name": illness.name,
                "Illness Start Date": new Date(illness.startDate).toLocaleDateString(),
                "Illness End Date": illness.endDate ? new Date(illness.endDate).toLocaleDateString() : "Ongoing",
                "Prescription UUID": "", "Prescription Issued At": "", "Prescription Illness Name": "", "Medications": ""
            };
            allData.push(Object.values(illnessRow));
        }

        for (const prescription of allPrescriptions) {
            const prescriptionRow = {
                "Record Type": "Prescription",
                "Patient First Name": patient.firstName,
                "Patient Last Name": patient.lastName,
                "Patient OIB": patient.oib,
                "Checkup UUID": "", "Checkup Date": "", "Checkup Type": "", "Checkup Illness ID": "",
                "Illness UUID": "", "Illness Name": "", "Illness Start Date": "", "Illness End Date": "",
                "Prescription UUID": prescription.uuid,
                "Prescription Issued At": new Date(prescription.issuedAt).toLocaleDateString(),
                "Prescription Illness Name": prescription.illnessName,
                "Medications": prescription.medications.map(m => m.name).join(', ')
            };
            allData.push(Object.values(prescriptionRow));
        }
    }

    const csvContent = allData.map(row => row.map(value => `"${value}"`).join(',')).join('\n');
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);
    link.setAttribute('href', url);
    link.setAttribute('download', 'all_patients_data.csv');
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
}