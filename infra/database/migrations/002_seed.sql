-- 002_seed.sql
-- Seed data: 3 Demo Users & 100+ Realistic Medicines

-- Demo Users
-- Passwords:
-- admin@pharmacy.local   => AdminPass123!
-- staff@pharmacy.local   => StaffPass123!
-- customer@pharmacy.local=> CustomerPass123!

INSERT INTO users (user_id, name, email, password_hash, role)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'Admin Supervisor', 'admin@pharmacy.local', '$2a$10$xCsG76MlUtBqvm8l/egWuerylO9AgdW7Hm7NWBtbVs3K44kgZ8eU6', 'admin'),
    ('22222222-2222-2222-2222-222222222222', 'Pharmacist Sarah Jenkins', 'staff@pharmacy.local', '$2a$10$sog4oSPNv.qpLTBb872DeuWypoxlkvRERk4rmFldpYSeCBUUZbbsy', 'staff'),
    ('33333333-3333-3333-3333-333333333333', 'John Doe (Customer)', 'customer@pharmacy.local', '$2a$10$.mFgHZ1EBw9x1YOhsyRzAuHTzu7bs0mftMKn8k74LtSU6h72U/mqC', 'customer')
ON CONFLICT (email) DO NOTHING;

-- Realistic Medicines (105 items)
INSERT INTO medicines (name, generic_name, description, price, stock_qty, requires_rx)
VALUES
    -- Analgesics & Anti-inflammatory (OTC & Rx)
    ('Paracetamol 500mg Tablet', 'Paracetamol', 'Analgesic and antipyretic oral tablet 500mg', 4.50, 250, FALSE),
    ('Paracetamol 1000mg Tablet', 'Paracetamol', 'Analgesic and antipyretic high-strength oral tablet 1000mg', 6.00, 180, FALSE),
    ('Ibuprofen 200mg Tablet', 'Ibuprofen', 'Non-steroidal anti-inflammatory drug (NSAID) 200mg tablet', 5.25, 200, FALSE),
    ('Ibuprofen 400mg Tablet', 'Ibuprofen', 'Non-steroidal anti-inflammatory drug (NSAID) 400mg tablet', 7.50, 150, FALSE),
    ('Ibuprofen 600mg Tablet', 'Ibuprofen', 'Prescription strength anti-inflammatory tablet 600mg', 12.00, 90, TRUE),
    ('Aspirin 81mg Tablet', 'Acetylsalicylic Acid', 'Low-dose enteric-coated tablet for cardiovascular prophylaxis', 4.00, 300, FALSE),
    ('Aspirin 325mg Tablet', 'Acetylsalicylic Acid', 'Analgesic and antiplatelet oral tablet 325mg', 5.00, 120, FALSE),
    ('Naproxen 220mg Tablet', 'Naproxen Sodium', 'Over-the-counter long-acting NSAID tablet 220mg', 8.50, 140, FALSE),
    ('Naproxen 500mg Tablet', 'Naproxen', 'Prescription anti-inflammatory oral tablet 500mg', 14.50, 75, TRUE),
    ('Celecoxib 200mg Capsule', 'Celecoxib', 'Selective COX-2 inhibitor anti-inflammatory capsule 200mg', 28.00, 60, TRUE),
    ('Diclofenac Sodium 50mg Tablet', 'Diclofenac Sodium', 'Enteric coated NSAID oral tablet 50mg', 11.20, 85, TRUE),
    ('Meloxicam 15mg Tablet', 'Meloxicam', 'Prescription oxicam NSAID tablet 15mg', 16.75, 70, TRUE),
    ('Tramadol 50mg Capsule', 'Tramadol Hydrochloride', 'Controlled opioid analgesic capsule 50mg', 22.50, 45, TRUE),
    ('Codeine Phosphate 30mg Tablet', 'Codeine Phosphate', 'Opioid analgesic tablet 30mg for moderate pain', 19.80, 40, TRUE),

    -- Antibiotics & Anti-Infectives (Rx)
    ('Amoxicillin 250mg Capsule', 'Amoxicillin', 'Broad-spectrum beta-lactam antibiotic capsule 250mg', 9.50, 110, TRUE),
    ('Amoxicillin 500mg Capsule', 'Amoxicillin', 'Broad-spectrum beta-lactam antibiotic capsule 500mg', 14.00, 150, TRUE),
    ('Amoxicillin-Clavulanate 875/125mg', 'Amoxicillin and Clavulanate Potassium', 'Potentiated penicillin antibiotic tablet', 24.50, 95, TRUE),
    ('Azithromycin 250mg Tablet', 'Azithromycin', 'Macrolide antibiotic film-coated tablet 250mg', 18.00, 80, TRUE),
    ('Azithromycin 500mg Tablet', 'Azithromycin', 'Macrolide antibiotic film-coated tablet 500mg', 26.50, 70, TRUE),
    ('Ciprofloxacin 500mg Tablet', 'Ciprofloxacin Hydrochloride', 'Fluoroquinolone antibiotic tablet 500mg', 17.50, 65, TRUE),
    ('Levofloxacin 500mg Tablet', 'Levofloxacin', 'Third-generation fluoroquinolone antibiotic tablet 500mg', 29.00, 50, TRUE),
    ('Cephalexin 500mg Capsule', 'Cephalexin', 'First-generation cephalosporin antibiotic capsule 500mg', 16.20, 85, TRUE),
    ('Cefuroxime Axetil 500mg Tablet', 'Cefuroxime Axetil', 'Second-generation cephalosporin antibiotic tablet 500mg', 32.00, 55, TRUE),
    ('Doxycycline Hyclate 100mg Capsule', 'Doxycycline', 'Tetracycline class broad-spectrum antibiotic capsule 100mg', 15.00, 100, TRUE),
    ('Metronidazole 400mg Tablet', 'Metronidazole', 'Nitroimidazole antimicrobial tablet 400mg', 8.90, 90, TRUE),
    ('Nitrofurantoin 100mg Capsule', 'Nitrofurantoin', 'Urinary tract anti-infective macrocrystal capsule 100mg', 21.00, 60, TRUE),
    ('Clarithromycin 500mg Tablet', 'Clarithromycin', 'Macrolide antibacterial tablet 500mg', 27.50, 45, TRUE),
    ('Fluconazole 150mg Capsule', 'Fluconazole', 'Systemic triazole antifungal capsule single dose 150mg', 12.00, 75, TRUE),

    -- Cardiovascular & Antihypertensive (Rx)
    ('Amlodipine 5mg Tablet', 'Amlodipine Besylate', 'Dihydropyridine calcium channel blocker tablet 5mg', 8.00, 200, TRUE),
    ('Amlodipine 10mg Tablet', 'Amlodipine Besylate', 'Dihydropyridine calcium channel blocker tablet 10mg', 12.50, 160, TRUE),
    ('Lisinopril 10mg Tablet', 'Lisinopril', 'ACE inhibitor antihypertensive tablet 10mg', 9.20, 180, TRUE),
    ('Lisinopril 20mg Tablet', 'Lisinopril', 'ACE inhibitor antihypertensive tablet 20mg', 13.40, 140, TRUE),
    ('Losartan Potassium 50mg Tablet', 'Losartan Potassium', 'Angiotensin II receptor antagonist (ARB) tablet 50mg', 11.50, 190, TRUE),
    ('Losartan Potassium 100mg Tablet', 'Losartan Potassium', 'Angiotensin II receptor antagonist (ARB) tablet 100mg', 17.00, 130, TRUE),
    ('Valsartan 80mg Tablet', 'Valsartan', 'Angiotensin II receptor antagonist tablet 80mg', 15.80, 80, TRUE),
    ('Telmisartan 40mg Tablet', 'Telmisartan', 'Long-acting angiotensin receptor blocker tablet 40mg', 18.50, 90, TRUE),
    ('Metoprolol Tartrate 50mg Tablet', 'Metoprolol Tartrate', 'Selective beta-1 adrenergic receptor blocker tablet 50mg', 10.00, 150, TRUE),
    ('Metoprolol Succinate ER 50mg', 'Metoprolol Succinate', 'Extended-release beta blocker tablet 50mg', 19.50, 110, TRUE),
    ('Atenolol 50mg Tablet', 'Atenolol', 'Beta-adrenergic antagonist tablet 50mg', 7.50, 120, TRUE),
    ('Hydrochlorothiazide 25mg Tablet', 'Hydrochlorothiazide', 'Thiazide diuretic antihypertensive tablet 25mg', 6.50, 175, TRUE),
    ('Furosemide 40mg Tablet', 'Furosemide', 'Loop diuretic oral tablet 40mg', 8.20, 130, TRUE),
    ('Spironolactone 25mg Tablet', 'Spironolactone', 'Potassium-sparing aldosterone antagonist tablet 25mg', 14.00, 85, TRUE),
    ('Atorvastatin 10mg Tablet', 'Atorvastatin Calcium', 'HMG-CoA reductase inhibitor lipid-lowering tablet 10mg', 11.00, 220, TRUE),
    ('Atorvastatin 20mg Tablet', 'Atorvastatin Calcium', 'HMG-CoA reductase inhibitor lipid-lowering tablet 20mg', 16.50, 200, TRUE),
    ('Atorvastatin 40mg Tablet', 'Atorvastatin Calcium', 'HMG-CoA reductase inhibitor lipid-lowering tablet 40mg', 23.00, 110, TRUE),
    ('Rosuvastatin 10mg Tablet', 'Rosuvastatin Calcium', 'Statin lipid-regulating tablet 10mg', 18.00, 140, TRUE),
    ('Rosuvastatin 20mg Tablet', 'Rosuvastatin Calcium', 'Statin lipid-regulating tablet 20mg', 25.50, 100, TRUE),
    ('Simvastatin 20mg Tablet', 'Simvastatin', 'Lipid-lowering agent tablet 20mg', 9.00, 130, TRUE),
    ('Clopidogrel 75mg Tablet', 'Clopidogrel Bisulfate', 'P2Y12 platelet inhibitor oral tablet 75mg', 22.00, 95, TRUE),

    -- Gastrointestinal (OTC & Rx)
    ('Omeprazole 20mg Capsule', 'Omeprazole', 'Proton pump inhibitor delayed-release capsule 20mg', 11.50, 180, FALSE),
    ('Omeprazole 40mg Capsule', 'Omeprazole', 'Proton pump inhibitor high-strength capsule 40mg', 18.00, 90, TRUE),
    ('Esomeprazole 20mg Tablet', 'Esomeprazole Magnesium', 'Delayed release acid reducer tablet 20mg', 14.50, 120, FALSE),
    ('Pantoprazole 40mg Tablet', 'Pantoprazole Sodium', 'Proton pump inhibitor delayed-release tablet 40mg', 16.00, 140, TRUE),
    ('Famotidine 20mg Tablet', 'Famotidine', 'Histamine H2 receptor antagonist tablet 20mg', 8.25, 160, FALSE),
    ('Ranitidine 150mg Tablet', 'Ranitidine Hydrochloride', 'H2 blocker oral tablet 150mg', 7.50, 80, FALSE),
    ('Antacid Chewable Tablet 500mg', 'Calcium Carbonate', 'Fast-acting antacid chewable fruit-flavored tablet', 4.50, 250, FALSE),
    ('Magnesium Hydroxide Suspension', 'Magnesium Hydroxide', 'Oral saline laxative and antacid liquid 350ml', 6.75, 75, FALSE),
    ('Loperamide 2mg Capsule', 'Loperamide Hydrochloride', 'Anti-diarrheal oral capsule 2mg', 5.50, 190, FALSE),
    ('Simethicone 80mg Chewable', 'Simethicone', 'Anti-gas antiflatulent chewable tablet 80mg', 5.00, 130, FALSE),
    ('Bisacodyl 5mg Tablet', 'Bisacodyl', 'Stimulant laxative enteric-coated tablet 5mg', 6.00, 110, FALSE),
    ('Ondansetron 4mg ODT', 'Ondansetron', 'Serotonin 5-HT3 antagonist orally disintegrating tablet', 28.50, 65, TRUE),

    -- Endocrine & Diabetes (Rx)
    ('Metformin 500mg Tablet', 'Metformin Hydrochloride', 'Biguanide oral antidiabetic tablet 500mg', 7.00, 250, TRUE),
    ('Metformin 850mg Tablet', 'Metformin Hydrochloride', 'Biguanide oral antidiabetic tablet 850mg', 10.50, 190, TRUE),
    ('Metformin ER 500mg Tablet', 'Metformin Hydrochloride', 'Extended-release oral antidiabetic tablet 500mg', 12.00, 160, TRUE),
    ('Glipizide 5mg Tablet', 'Glipizide', 'Sulfonylurea oral hypoglycemic tablet 5mg', 9.50, 110, TRUE),
    ('Sitagliptin 100mg Tablet', 'Sitagliptin Phosphate', 'DPP-4 inhibitor antidiabetic tablet 100mg', 45.00, 80, TRUE),
    ('Levothyroxine 25mcg Tablet', 'Levothyroxine Sodium', 'Synthetic thyroid hormone T4 tablet 25mcg', 8.00, 120, TRUE),
    ('Levothyroxine 50mcg Tablet', 'Levothyroxine Sodium', 'Synthetic thyroid hormone T4 tablet 50mcg', 10.00, 180, TRUE),
    ('Levothyroxine 100mcg Tablet', 'Levothyroxine Sodium', 'Synthetic thyroid hormone T4 tablet 100mcg', 12.50, 150, TRUE),

    -- Respiratory & Allergy (OTC & Rx)
    ('Cetirizine 10mg Tablet', 'Cetirizine Hydrochloride', 'Second-generation non-drowsy antihistamine tablet 10mg', 7.00, 220, FALSE),
    ('Loratadine 10mg Tablet', 'Loratadine', 'Non-sedating antihistamine tablet 10mg', 6.50, 210, FALSE),
    ('Fexofenadine 180mg Tablet', 'Fexofenadine Hydrochloride', 'Non-sedating 24hr allergy relief tablet 180mg', 14.00, 130, FALSE),
    ('Diphenhydramine 25mg Capsule', 'Diphenhydramine Hydrochloride', 'First-generation antihistamine capsule 25mg', 4.75, 175, FALSE),
    ('Montelukast 10mg Tablet', 'Montelukast Sodium', 'Leukotriene receptor antagonist oral tablet 10mg', 19.50, 105, TRUE),
    ('Fluticasone Propionate Nasal Spray', 'Fluticasone Propionate', 'Glucocorticoid nasal allergy spray 50mcg per spray', 16.50, 95, FALSE),
    ('Salbutamol Inhaler 100mcg', 'Salbutamol (Albuterol)', 'Short-acting beta-2 agonist bronchodilator inhaler 200 doses', 21.00, 115, TRUE),
    ('Budesonide Turbuhaler 200mcg', 'Budesonide', 'Inhaled corticosteroid maintenance inhaler 200mcg', 38.00, 50, TRUE),
    ('Guaifenesin 400mg Tablet', 'Guaifenesin', 'Expectorant mucus relief oral tablet 400mg', 6.00, 140, FALSE),
    ('Dextromethorphan Cough Syrup 150ml', 'Dextromethorphan HBr', 'Oral antitussive cough suppressant liquid 15mg/5ml', 8.50, 85, FALSE),
    ('Oxymetazoline Nasal Spray 15ml', 'Oxymetazoline Hydrochloride', 'Topical nasal decongestant spray 0.05%', 7.20, 95, FALSE),

    -- Dermatology & Topical (OTC & Rx)
    ('Hydrocortisone 1% Cream 30g', 'Hydrocortisone', 'Mild topical corticosteroid anti-itch cream 1%', 6.50, 120, FALSE),
    ('Betamethasone 0.05% Cream 15g', 'Betamethasone Dipropionate', 'Potent topical corticosteroid cream 0.05%', 15.00, 60, TRUE),
    ('Clotrimazole 1% Cream 20g', 'Clotrimazole', 'Broad-spectrum topical antifungal cream 1%', 5.50, 130, FALSE),
    ('Terbinafine 1% Cream 15g', 'Terbinafine Hydrochloride', 'Topical allylamine antifungal cream 1%', 9.75, 80, FALSE),
    ('Mupirocin 2% Ointment 15g', 'Mupirocin', 'Topical antibacterial ointment 2%', 17.50, 70, TRUE),
    ('Permethrin 5% Dermal Cream 30g', 'Permethrin', 'Scabicidal topical cream 5%', 22.00, 40, TRUE),
    ('Calamine Lotion 100ml', 'Calamine and Zinc Oxide', 'Soothing topical protective anti-itch suspension', 4.00, 110, FALSE),

    -- Ophthalmology & Otic (OTC & Rx)
    ('Carboxymethylcellulose 0.5% Drops', 'Carboxymethylcellulose Sodium', 'Lubricating artificial tears eye drops 15ml', 7.50, 140, FALSE),
    ('Olopatadine 0.1% Eye Drops 5ml', 'Olopatadine Hydrochloride', 'Ophthalmic antihistamine allergy drops 0.1%', 24.00, 55, TRUE),
    ('Tobramycin 0.3% Eye Drops 5ml', 'Tobramycin', 'Ophthalmic aminoglycoside antibacterial drops 0.3%', 16.50, 65, TRUE),
    ('Ciprofloxacin 0.3% Otic Drops 10ml', 'Ciprofloxacin', 'Antibacterial ear solution 0.3%', 18.00, 50, TRUE),

    -- Central Nervous System & Psychiatric (Rx)
    ('Sertraline 50mg Tablet', 'Sertraline Hydrochloride', 'Selective serotonin reuptake inhibitor (SSRI) tablet 50mg', 15.00, 120, TRUE),
    ('Escitalopram 10mg Tablet', 'Escitalopram Oxalate', 'Selective serotonin reuptake inhibitor (SSRI) tablet 10mg', 18.20, 100, TRUE),
    ('Fluoxetine 20mg Capsule', 'Fluoxetine Hydrochloride', 'Selective serotonin reuptake inhibitor capsule 20mg', 14.50, 90, TRUE),
    ('Duloxetine 30mg Capsule', 'Duloxetine Hydrochloride', 'Serotonin-norepinephrine reuptake inhibitor capsule 30mg', 27.00, 60, TRUE),
    ('Gabapentin 300mg Capsule', 'Gabapentin', 'Anticonvulsant and neuropathic pain agent capsule 300mg', 19.50, 85, TRUE),
    ('Pregabalin 75mg Capsule', 'Pregabalin', 'Anticonvulsant and anxiolytic capsule 75mg', 34.00, 55, TRUE),
    ('Zolpidem 10mg Tablet', 'Zolpidem Tartrate', 'Controlled non-benzodiazepine sedative hypnotic tablet 10mg', 26.50, 40, TRUE),

    -- Vitamins, Minerals & Nutrition (OTC)
    ('Vitamin C 500mg Tablet', 'Ascorbic Acid', 'Antioxidant immune supplement tablet 500mg', 5.00, 200, FALSE),
    ('Vitamin D3 1000IU Capsule', 'Cholecalciferol', 'Bone health dietary supplement softgel 1000IU', 7.25, 230, FALSE),
    ('Vitamin D3 50000IU Capsule', 'Cholecalciferol', 'High-potency weekly vitamin D3 supplement capsule', 18.00, 75, TRUE),
    ('Vitamin B-Complex Tablet', 'B-Complex Vitamins', 'Daily nutritional metabolic support tablets', 6.50, 180, FALSE),
    ('Zinc Sulfate 50mg Tablet', 'Zinc Sulfate', 'Dietary mineral supplement tablet 50mg', 4.75, 160, FALSE),
    ('Ferrous Sulfate 325mg Tablet', 'Ferrous Sulfate', 'Oral iron deficiency supplement tablet 325mg (65mg elemental)', 6.20, 140, FALSE),
    ('Folic Acid 5mg Tablet', 'Folic Acid', 'Essential folate supplement tablet 5mg', 5.00, 150, FALSE),
    ('Calcium + Vitamin D3 Tablet', 'Calcium Carbonate with Vitamin D3', 'Bone density dietary supplement tablet', 8.50, 170, FALSE),
    ('Multivitamin Daily Formula', 'Multivitamins & Minerals', 'Comprehensive daily oral dietary supplement caplets', 9.90, 190, FALSE)
ON CONFLICT DO NOTHING;
