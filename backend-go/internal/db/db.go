package db

import (
	"context"
	"fmt"
	"log"
	"smart-pharmacy/backend-go/internal/models"
	"time"
)

// InitDB initializes PostgreSQL connection pool or returns an error
func InitDB(ctx context.Context, databaseURL string) (Store, error) {
	if databaseURL == "" || databaseURL == "memory" {
		log.Println("[DB] Using thread-safe in-memory store with pre-seeded medicines and demo users")
		mem := NewMemoryStore()
		SeedMemoryStore(mem)
		return mem, nil
	}

	pgStore, err := NewPostgresStore(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to PostgreSQL at %s: %w", databaseURL, err)
	}

	log.Println("[DB] Successfully connected to PostgreSQL")
	return pgStore, nil
}

// SeedMemoryStore seeds realistic medicines and demo users into an in-memory store
func SeedMemoryStore(m *MemoryStore) {
	ctx := context.Background()

	// Demo Users (hashes match AdminPass123!, StaffPass123!, CustomerPass123!)
	users := []models.User{
		{
			UserID:       "11111111-1111-1111-1111-111111111111",
			Name:         "Admin Supervisor",
			Email:        "admin@pharmacy.local",
			PasswordHash: "$2a$10$xCsG76MlUtBqvm8l/egWuerylO9AgdW7Hm7NWBtbVs3K44kgZ8eU6",
			Role:         models.RoleAdmin,
			CreatedAt:    time.Now().UTC(),
		},
		{
			UserID:       "22222222-2222-2222-2222-222222222222",
			Name:         "Pharmacist Sarah Jenkins",
			Email:        "staff@pharmacy.local",
			PasswordHash: "$2a$10$sog4oSPNv.qpLTBb872DeuWypoxlkvRERk4rmFldpYSeCBUUZbbsy",
			Role:         models.RoleStaff,
			CreatedAt:    time.Now().UTC(),
		},
		{
			UserID:       "33333333-3333-3333-3333-333333333333",
			Name:         "John Doe (Customer)",
			Email:        "customer@pharmacy.local",
			PasswordHash: "$2a$10$.mFgHZ1EBw9x1YOhsyRzAuHTzu7bs0mftMKn8k74LtSU6h72U/mqC",
			Role:         models.RoleCustomer,
			CreatedAt:    time.Now().UTC(),
		},
	}

	for _, u := range users {
		_ = m.CreateUser(ctx, &u)
	}

	// 100+ realistic medicines
	medicines := []models.Medicine{
		{Name: "Paracetamol 500mg Tablet", GenericName: "Paracetamol", Description: "Analgesic and antipyretic oral tablet 500mg", Price: 4.50, StockQty: 250, RequiresRx: false},
		{Name: "Paracetamol 1000mg Tablet", GenericName: "Paracetamol", Description: "Analgesic and antipyretic high-strength oral tablet 1000mg", Price: 6.00, StockQty: 180, RequiresRx: false},
		{Name: "Ibuprofen 200mg Tablet", GenericName: "Ibuprofen", Description: "Non-steroidal anti-inflammatory drug (NSAID) 200mg tablet", Price: 5.25, StockQty: 200, RequiresRx: false},
		{Name: "Ibuprofen 400mg Tablet", GenericName: "Ibuprofen", Description: "Non-steroidal anti-inflammatory drug (NSAID) 400mg tablet", Price: 7.50, StockQty: 150, RequiresRx: false},
		{Name: "Ibuprofen 600mg Tablet", GenericName: "Ibuprofen", Description: "Prescription strength anti-inflammatory tablet 600mg", Price: 12.00, StockQty: 90, RequiresRx: true},
		{Name: "Aspirin 81mg Tablet", GenericName: "Acetylsalicylic Acid", Description: "Low-dose enteric-coated tablet for cardiovascular prophylaxis", Price: 4.00, StockQty: 300, RequiresRx: false},
		{Name: "Aspirin 325mg Tablet", GenericName: "Acetylsalicylic Acid", Description: "Analgesic and antiplatelet oral tablet 325mg", Price: 5.00, StockQty: 120, RequiresRx: false},
		{Name: "Naproxen 220mg Tablet", GenericName: "Naproxen Sodium", Description: "Over-the-counter long-acting NSAID tablet 220mg", Price: 8.50, StockQty: 140, RequiresRx: false},
		{Name: "Naproxen 500mg Tablet", GenericName: "Naproxen", Description: "Prescription anti-inflammatory oral tablet 500mg", Price: 14.50, StockQty: 75, RequiresRx: true},
		{Name: "Celecoxib 200mg Capsule", GenericName: "Celecoxib", Description: "Selective COX-2 inhibitor anti-inflammatory capsule 200mg", Price: 28.00, StockQty: 60, RequiresRx: true},
		{Name: "Diclofenac Sodium 50mg Tablet", GenericName: "Diclofenac Sodium", Description: "Enteric coated NSAID oral tablet 50mg", Price: 11.20, StockQty: 85, RequiresRx: true},
		{Name: "Meloxicam 15mg Tablet", GenericName: "Meloxicam", Description: "Prescription oxicam NSAID tablet 15mg", Price: 16.75, StockQty: 70, RequiresRx: true},
		{Name: "Tramadol 50mg Capsule", GenericName: "Tramadol Hydrochloride", Description: "Controlled opioid analgesic capsule 50mg", Price: 22.50, StockQty: 45, RequiresRx: true},
		{Name: "Codeine Phosphate 30mg Tablet", GenericName: "Codeine Phosphate", Description: "Opioid analgesic tablet 30mg for moderate pain", Price: 19.80, StockQty: 40, RequiresRx: true},

		{Name: "Amoxicillin 250mg Capsule", GenericName: "Amoxicillin", Description: "Broad-spectrum beta-lactam antibiotic capsule 250mg", Price: 9.50, StockQty: 110, RequiresRx: true},
		{Name: "Amoxicillin 500mg Capsule", GenericName: "Amoxicillin", Description: "Broad-spectrum beta-lactam antibiotic capsule 500mg", Price: 14.00, StockQty: 150, RequiresRx: true},
		{Name: "Amoxicillin-Clavulanate 875/125mg", GenericName: "Amoxicillin and Clavulanate Potassium", Description: "Potentiated penicillin antibiotic tablet", Price: 24.50, StockQty: 95, RequiresRx: true},
		{Name: "Azithromycin 250mg Tablet", GenericName: "Azithromycin", Description: "Macrolide antibiotic film-coated tablet 250mg", Price: 18.00, StockQty: 80, RequiresRx: true},
		{Name: "Azithromycin 500mg Tablet", GenericName: "Azithromycin", Description: "Macrolide antibiotic film-coated tablet 500mg", Price: 26.50, StockQty: 70, RequiresRx: true},
		{Name: "Ciprofloxacin 500mg Tablet", GenericName: "Ciprofloxacin Hydrochloride", Description: "Fluoroquinolone antibiotic tablet 500mg", Price: 17.50, StockQty: 65, RequiresRx: true},
		{Name: "Levofloxacin 500mg Tablet", GenericName: "Levofloxacin", Description: "Third-generation fluoroquinolone antibiotic tablet 500mg", Price: 29.00, StockQty: 50, RequiresRx: true},
		{Name: "Cephalexin 500mg Capsule", GenericName: "Cephalexin", Description: "First-generation cephalosporin antibiotic capsule 500mg", Price: 16.20, StockQty: 85, RequiresRx: true},
		{Name: "Cefuroxime Axetil 500mg Tablet", GenericName: "Cefuroxime Axetil", Description: "Second-generation cephalosporin antibiotic tablet 500mg", Price: 32.00, StockQty: 55, RequiresRx: true},
		{Name: "Doxycycline Hyclate 100mg Capsule", GenericName: "Doxycycline", Description: "Tetracycline class broad-spectrum antibiotic capsule 100mg", Price: 15.00, StockQty: 100, RequiresRx: true},
		{Name: "Metronidazole 400mg Tablet", GenericName: "Metronidazole", Description: "Nitroimidazole antimicrobial tablet 400mg", Price: 8.90, StockQty: 90, RequiresRx: true},
		{Name: "Nitrofurantoin 100mg Capsule", GenericName: "Nitrofurantoin", Description: "Urinary tract anti-infective macrocrystal capsule 100mg", Price: 21.00, StockQty: 60, RequiresRx: true},
		{Name: "Clarithromycin 500mg Tablet", GenericName: "Clarithromycin", Description: "Macrolide antibacterial tablet 500mg", Price: 27.50, StockQty: 45, RequiresRx: true},
		{Name: "Fluconazole 150mg Capsule", GenericName: "Fluconazole", Description: "Systemic triazole antifungal capsule single dose 150mg", Price: 12.00, StockQty: 75, RequiresRx: true},

		{Name: "Amlodipine 5mg Tablet", GenericName: "Amlodipine Besylate", Description: "Dihydropyridine calcium channel blocker tablet 5mg", Price: 8.00, StockQty: 200, RequiresRx: true},
		{Name: "Amlodipine 10mg Tablet", GenericName: "Amlodipine Besylate", Description: "Dihydropyridine calcium channel blocker tablet 10mg", Price: 12.50, StockQty: 160, RequiresRx: true},
		{Name: "Lisinopril 10mg Tablet", GenericName: "Lisinopril", Description: "ACE inhibitor antihypertensive tablet 10mg", Price: 9.20, StockQty: 180, RequiresRx: true},
		{Name: "Lisinopril 20mg Tablet", GenericName: "Lisinopril", Description: "ACE inhibitor antihypertensive tablet 20mg", Price: 13.40, StockQty: 140, RequiresRx: true},
		{Name: "Losartan Potassium 50mg Tablet", GenericName: "Losartan Potassium", Description: "Angiotensin II receptor antagonist (ARB) tablet 50mg", Price: 11.50, StockQty: 190, RequiresRx: true},
		{Name: "Losartan Potassium 100mg Tablet", GenericName: "Losartan Potassium", Description: "Angiotensin II receptor antagonist (ARB) tablet 100mg", Price: 17.00, StockQty: 130, RequiresRx: true},
		{Name: "Valsartan 80mg Tablet", GenericName: "Valsartan", Description: "Angiotensin II receptor antagonist tablet 80mg", Price: 15.80, StockQty: 80, RequiresRx: true},
		{Name: "Telmisartan 40mg Tablet", GenericName: "Telmisartan", Description: "Long-acting angiotensin receptor blocker tablet 40mg", Price: 18.50, StockQty: 90, RequiresRx: true},
		{Name: "Metoprolol Tartrate 50mg Tablet", GenericName: "Metoprolol Tartrate", Description: "Selective beta-1 adrenergic receptor blocker tablet 50mg", Price: 10.00, StockQty: 150, RequiresRx: true},
		{Name: "Metoprolol Succinate ER 50mg", GenericName: "Metoprolol Succinate", Description: "Extended-release beta blocker tablet 50mg", Price: 19.50, StockQty: 110, RequiresRx: true},
		{Name: "Atenolol 50mg Tablet", GenericName: "Atenolol", Description: "Beta-adrenergic antagonist tablet 50mg", Price: 7.50, StockQty: 120, RequiresRx: true},
		{Name: "Hydrochlorothiazide 25mg Tablet", GenericName: "Hydrochlorothiazide", Description: "Thiazide diuretic antihypertensive tablet 25mg", Price: 6.50, StockQty: 175, RequiresRx: true},
		{Name: "Furosemide 40mg Tablet", GenericName: "Furosemide", Description: "Loop diuretic oral tablet 40mg", Price: 8.20, StockQty: 130, RequiresRx: true},
		{Name: "Spironolactone 25mg Tablet", GenericName: "Spironolactone", Description: "Potassium-sparing aldosterone antagonist tablet 25mg", Price: 14.00, StockQty: 85, RequiresRx: true},
		{Name: "Atorvastatin 10mg Tablet", GenericName: "Atorvastatin Calcium", Description: "HMG-CoA reductase inhibitor lipid-lowering tablet 10mg", Price: 11.00, StockQty: 220, RequiresRx: true},
		{Name: "Atorvastatin 20mg Tablet", GenericName: "Atorvastatin Calcium", Description: "HMG-CoA reductase inhibitor lipid-lowering tablet 20mg", Price: 16.50, StockQty: 200, RequiresRx: true},
		{Name: "Atorvastatin 40mg Tablet", GenericName: "Atorvastatin Calcium", Description: "HMG-CoA reductase inhibitor lipid-lowering tablet 40mg", Price: 23.00, StockQty: 110, RequiresRx: true},
		{Name: "Rosuvastatin 10mg Tablet", GenericName: "Rosuvastatin Calcium", Description: "Statin lipid-regulating tablet 10mg", Price: 18.00, StockQty: 140, RequiresRx: true},
		{Name: "Rosuvastatin 20mg Tablet", GenericName: "Rosuvastatin Calcium", Description: "Statin lipid-regulating tablet 20mg", Price: 25.50, StockQty: 100, RequiresRx: true},
		{Name: "Simvastatin 20mg Tablet", GenericName: "Simvastatin", Description: "Lipid-lowering agent tablet 20mg", Price: 9.00, StockQty: 130, RequiresRx: true},
		{Name: "Clopidogrel 75mg Tablet", GenericName: "Clopidogrel Bisulfate", Description: "P2Y12 platelet inhibitor oral tablet 75mg", Price: 22.00, StockQty: 95, RequiresRx: true},

		{Name: "Omeprazole 20mg Capsule", GenericName: "Omeprazole", Description: "Proton pump inhibitor delayed-release capsule 20mg", Price: 11.50, StockQty: 180, RequiresRx: false},
		{Name: "Omeprazole 40mg Capsule", GenericName: "Omeprazole", Description: "Proton pump inhibitor high-strength capsule 40mg", Price: 18.00, StockQty: 90, RequiresRx: true},
		{Name: "Esomeprazole 20mg Tablet", GenericName: "Esomeprazole Magnesium", Description: "Delayed release acid reducer tablet 20mg", Price: 14.50, StockQty: 120, RequiresRx: false},
		{Name: "Pantoprazole 40mg Tablet", GenericName: "Pantoprazole Sodium", Description: "Proton pump inhibitor delayed-release tablet 40mg", Price: 16.00, StockQty: 140, RequiresRx: true},
		{Name: "Famotidine 20mg Tablet", GenericName: "Famotidine", Description: "Histamine H2 receptor antagonist tablet 20mg", Price: 8.25, StockQty: 160, RequiresRx: false},
		{Name: "Ranitidine 150mg Tablet", GenericName: "Ranitidine Hydrochloride", Description: "H2 blocker oral tablet 150mg", Price: 7.50, StockQty: 80, RequiresRx: false},
		{Name: "Antacid Chewable Tablet 500mg", GenericName: "Calcium Carbonate", Description: "Fast-acting antacid chewable fruit-flavored tablet", Price: 4.50, StockQty: 250, RequiresRx: false},
		{Name: "Magnesium Hydroxide Suspension", GenericName: "Magnesium Hydroxide", Description: "Oral saline laxative and antacid liquid 350ml", Price: 6.75, StockQty: 75, RequiresRx: false},
		{Name: "Loperamide 2mg Capsule", GenericName: "Loperamide Hydrochloride", Description: "Anti-diarrheal oral capsule 2mg", Price: 5.50, StockQty: 190, RequiresRx: false},
		{Name: "Simethicone 80mg Chewable", GenericName: "Simethicone", Description: "Anti-gas antiflatulent chewable tablet 80mg", Price: 5.00, StockQty: 130, RequiresRx: false},
		{Name: "Bisacodyl 5mg Tablet", GenericName: "Bisacodyl", Description: "Stimulant laxative enteric-coated tablet 5mg", Price: 6.00, StockQty: 110, RequiresRx: false},
		{Name: "Ondansetron 4mg ODT", GenericName: "Ondansetron", Description: "Serotonin 5-HT3 antagonist orally disintegrating tablet", Price: 28.50, StockQty: 65, RequiresRx: true},

		{Name: "Metformin 500mg Tablet", GenericName: "Metformin Hydrochloride", Description: "Biguanide oral antidiabetic tablet 500mg", Price: 7.00, StockQty: 250, RequiresRx: true},
		{Name: "Metformin 850mg Tablet", GenericName: "Metformin Hydrochloride", Description: "Biguanide oral antidiabetic tablet 850mg", Price: 10.50, StockQty: 190, RequiresRx: true},
		{Name: "Metformin ER 500mg Tablet", GenericName: "Metformin Hydrochloride", Description: "Extended-release oral antidiabetic tablet 500mg", Price: 12.00, StockQty: 160, RequiresRx: true},
		{Name: "Glipizide 5mg Tablet", GenericName: "Glipizide", Description: "Sulfonylurea oral hypoglycemic tablet 5mg", Price: 9.50, StockQty: 110, RequiresRx: true},
		{Name: "Sitagliptin 100mg Tablet", GenericName: "Sitagliptin Phosphate", Description: "DPP-4 inhibitor antidiabetic tablet 100mg", Price: 45.00, StockQty: 80, RequiresRx: true},
		{Name: "Levothyroxine 25mcg Tablet", GenericName: "Levothyroxine Sodium", Description: "Synthetic thyroid hormone T4 tablet 25mcg", Price: 8.00, StockQty: 120, RequiresRx: true},
		{Name: "Levothyroxine 50mcg Tablet", GenericName: "Levothyroxine Sodium", Description: "Synthetic thyroid hormone T4 tablet 50mcg", Price: 10.00, StockQty: 180, RequiresRx: true},
		{Name: "Levothyroxine 100mcg Tablet", GenericName: "Levothyroxine Sodium", Description: "Synthetic thyroid hormone T4 tablet 100mcg", Price: 12.50, StockQty: 150, RequiresRx: true},

		{Name: "Cetirizine 10mg Tablet", GenericName: "Cetirizine Hydrochloride", Description: "Second-generation non-drowsy antihistamine tablet 10mg", Price: 7.00, StockQty: 220, RequiresRx: false},
		{Name: "Loratadine 10mg Tablet", GenericName: "Loratadine", Description: "Non-sedating antihistamine tablet 10mg", Price: 6.50, StockQty: 210, RequiresRx: false},
		{Name: "Fexofenadine 180mg Tablet", GenericName: "Fexofenadine Hydrochloride", Description: "Non-sedating 24hr allergy relief tablet 180mg", Price: 14.00, StockQty: 130, RequiresRx: false},
		{Name: "Diphenhydramine 25mg Capsule", GenericName: "Diphenhydramine Hydrochloride", Description: "First-generation antihistamine capsule 25mg", Price: 4.75, StockQty: 175, RequiresRx: false},
		{Name: "Montelukast 10mg Tablet", GenericName: "Montelukast Sodium", Description: "Leukotriene receptor antagonist oral tablet 10mg", Price: 19.50, StockQty: 105, RequiresRx: true},
		{Name: "Fluticasone Propionate Nasal Spray", GenericName: "Fluticasone Propionate", Description: "Glucocorticoid nasal allergy spray 50mcg per spray", Price: 16.50, StockQty: 95, RequiresRx: false},
		{Name: "Salbutamol Inhaler 100mcg", GenericName: "Salbutamol (Albuterol)", Description: "Short-acting beta-2 agonist bronchodilator inhaler 200 doses", Price: 21.00, StockQty: 115, RequiresRx: true},
		{Name: "Budesonide Turbuhaler 200mcg", GenericName: "Budesonide", Description: "Inhaled corticosteroid maintenance inhaler 200mcg", Price: 38.00, StockQty: 50, RequiresRx: true},
		{Name: "Guaifenesin 400mg Tablet", GenericName: "Guaifenesin", Description: "Expectorant mucus relief oral tablet 400mg", Price: 6.00, StockQty: 140, RequiresRx: false},
		{Name: "Dextromethorphan Cough Syrup 150ml", GenericName: "Dextromethorphan HBr", Description: "Oral antitussive cough suppressant liquid 15mg/5ml", Price: 8.50, StockQty: 85, RequiresRx: false},
		{Name: "Oxymetazoline Nasal Spray 15ml", GenericName: "Oxymetazoline Hydrochloride", Description: "Topical nasal decongestant spray 0.05%", Price: 7.20, StockQty: 95, RequiresRx: false},

		{Name: "Hydrocortisone 1% Cream 30g", GenericName: "Hydrocortisone", Description: "Mild topical corticosteroid anti-itch cream 1%", Price: 6.50, StockQty: 120, RequiresRx: false},
		{Name: "Betamethasone 0.05% Cream 15g", GenericName: "Betamethasone Dipropionate", Description: "Potent topical corticosteroid cream 0.05%", Price: 15.00, StockQty: 60, RequiresRx: true},
		{Name: "Clotrimazole 1% Cream 20g", GenericName: "Clotrimazole", Description: "Broad-spectrum topical antifungal cream 1%", Price: 5.50, StockQty: 130, RequiresRx: false},
		{Name: "Terbinafine 1% Cream 15g", GenericName: "Terbinafine Hydrochloride", Description: "Topical allylamine antifungal cream 1%", Price: 9.75, StockQty: 80, RequiresRx: false},
		{Name: "Mupirocin 2% Ointment 15g", GenericName: "Mupirocin", Description: "Topical antibacterial ointment 2%", Price: 17.50, StockQty: 70, RequiresRx: true},
		{Name: "Permethrin 5% Dermal Cream 30g", GenericName: "Permethrin", Description: "Scabicidal topical cream 5%", Price: 22.00, StockQty: 40, RequiresRx: true},
		{Name: "Calamine Lotion 100ml", GenericName: "Calamine and Zinc Oxide", Description: "Soothing topical protective anti-itch suspension", Price: 4.00, StockQty: 110, RequiresRx: false},

		{Name: "Carboxymethylcellulose 0.5% Drops", GenericName: "Carboxymethylcellulose Sodium", Description: "Lubricating artificial tears eye drops 15ml", Price: 7.50, StockQty: 140, RequiresRx: false},
		{Name: "Olopatadine 0.1% Eye Drops 5ml", GenericName: "Olopatadine Hydrochloride", Description: "Ophthalmic antihistamine allergy drops 0.1%", Price: 24.00, StockQty: 55, RequiresRx: true},
		{Name: "Tobramycin 0.3% Eye Drops 5ml", GenericName: "Tobramycin", Description: "Ophthalmic aminoglycoside antibacterial drops 0.3%", Price: 16.50, StockQty: 65, RequiresRx: true},
		{Name: "Ciprofloxacin 0.3% Otic Drops 10ml", GenericName: "Ciprofloxacin", Description: "Antibacterial ear solution 0.3%", Price: 18.00, StockQty: 50, RequiresRx: true},

		{Name: "Sertraline 50mg Tablet", GenericName: "Sertraline Hydrochloride", Description: "Selective serotonin reuptake inhibitor (SSRI) tablet 50mg", Price: 15.00, StockQty: 120, RequiresRx: true},
		{Name: "Escitalopram 10mg Tablet", GenericName: "Escitalopram Oxalate", Description: "Selective serotonin reuptake inhibitor (SSRI) tablet 10mg", Price: 18.20, StockQty: 100, RequiresRx: true},
		{Name: "Fluoxetine 20mg Capsule", GenericName: "Fluoxetine Hydrochloride", Description: "Selective serotonin reuptake inhibitor capsule 20mg", Price: 14.50, StockQty: 90, RequiresRx: true},
		{Name: "Duloxetine 30mg Capsule", GenericName: "Duloxetine Hydrochloride", Description: "Serotonin-norepinephrine reuptake inhibitor capsule 30mg", Price: 27.00, StockQty: 60, RequiresRx: true},
		{Name: "Gabapentin 300mg Capsule", GenericName: "Gabapentin", Description: "Anticonvulsant and neuropathic pain agent capsule 300mg", Price: 19.50, StockQty: 85, RequiresRx: true},
		{Name: "Pregabalin 75mg Capsule", GenericName: "Pregabalin", Description: "Anticonvulsant and anxiolytic capsule 75mg", Price: 34.00, StockQty: 55, RequiresRx: true},
		{Name: "Zolpidem 10mg Tablet", GenericName: "Zolpidem Tartrate", Description: "Controlled non-benzodiazepine sedative hypnotic tablet 10mg", Price: 26.50, StockQty: 40, RequiresRx: true},

		{Name: "Vitamin C 500mg Tablet", GenericName: "Ascorbic Acid", Description: "Antioxidant immune supplement tablet 500mg", Price: 5.00, StockQty: 200, RequiresRx: false},
		{Name: "Vitamin D3 1000IU Capsule", GenericName: "Cholecalciferol", Description: "Bone health dietary supplement softgel 1000IU", Price: 7.25, StockQty: 230, RequiresRx: false},
		{Name: "Vitamin D3 50000IU Capsule", GenericName: "Cholecalciferol", Description: "High-potency weekly vitamin D3 supplement capsule", Price: 18.00, StockQty: 75, RequiresRx: true},
		{Name: "Vitamin B-Complex Tablet", GenericName: "B-Complex Vitamins", Description: "Daily nutritional metabolic support tablets", Price: 6.50, StockQty: 180, RequiresRx: false},
		{Name: "Zinc Sulfate 50mg Tablet", GenericName: "Zinc Sulfate", Description: "Dietary mineral supplement tablet 50mg", Price: 4.75, StockQty: 160, RequiresRx: false},
		{Name: "Ferrous Sulfate 325mg Tablet", GenericName: "Ferrous Sulfate", Description: "Oral iron deficiency supplement tablet 325mg (65mg elemental)", Price: 6.20, StockQty: 140, RequiresRx: false},
		{Name: "Folic Acid 5mg Tablet", GenericName: "Folic Acid", Description: "Essential folate supplement tablet 5mg", Price: 5.00, StockQty: 150, RequiresRx: false},
		{Name: "Calcium + Vitamin D3 Tablet", GenericName: "Calcium Carbonate with Vitamin D3", Description: "Bone density dietary supplement tablet", Price: 8.50, StockQty: 170, RequiresRx: false},
		{Name: "Multivitamin Daily Formula", GenericName: "Multivitamins & Minerals", Description: "Comprehensive daily oral dietary supplement caplets", Price: 9.90, StockQty: 190, RequiresRx: false},
	}

	for _, med := range medicines {
		_ = m.CreateMedicine(ctx, &med)
	}
}
