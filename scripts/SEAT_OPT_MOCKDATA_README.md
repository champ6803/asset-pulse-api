# Seat Optimization Mock Data

ไฟล์ SQL สำหรับสร้าง mock data จำนวนมากเพื่อทดสอบหน้า Seat Optimization

## ไฟล์ที่เกี่ยวข้อง

1. **0004_insert_seat_opt_mockdata.sql** - Basic seat optimization mock data
2. **0005_insert_massive_seat_opt.sql** - Massive seat optimization mock data (แนะนำให้ใช้ตัวนี้)

## วิธีใช้งาน

### 1. รัน SQL Scripts ตามลำดับ

```bash
# ไปที่โฟลเดอร์ scripts
cd asset-pulse-api/scripts

# เชื่อมต่อ database และรัน script
psql -U postgres -d assetpulse -f 0001_initial_tables.sql
psql -U postgres -d assetpulse -f 0002_insert_mockdata.sql
psql -U postgres -d assetpulse -f 0003_create_view_recommend_pack_for_user.sql
psql -U postgres -d assetpulse -f 0005_insert_massive_seat_opt.sql  # ⭐ แนะนำตัวนี้ (รองรับการรันซ้ำ)
```

### ⚠️ หมายเหตุ
- Scripts มีการลบข้อมูลเก่า (cleanup) ก่อน insert ใหม่ **รองรับการรันซ้ำได้**
- ถ้าจะรันซ้ำก็ใช้คำสั่งเดียวได้เลย ไม่มี duplicate key error

### 2. หรือรันผ่าน Docker Compose

```bash
# From project root
cd asset-pulse-api
docker-compose exec postgres psql -U postgres -d assetpulse -f /scripts/0005_insert_massive_seat_opt.sql
```

## สิ่งที่ Script จะสร้าง

### 0004_insert_seat_opt_mockdata.sql
- สร้าง inactive users (5 คนต่อ department)
- สร้าง license assignments สำหรับ inactive users
- สร้าง reallocate opportunities (ระหว่าง departments)
- สร้าง downgrade opportunities (Pro/Enterprise -> lower tier)
- เพิ่มราคา license tiers (เพื่อให้เห็น savings เยอะขึ้น)

### 0005_insert_massive_seat_opt.sql (แนะนำ)
- **Revoke Opportunities**: Users ที่มี licenses แต่ไม่ได้ใช้ 90+ วัน
  - จำนวน: ~200-300+ opportunities
  - เงินประหยัด: ~3000-4000 THB/user/month สำหรับ Enterprise
  - เงินประหยัด: ~1500-2000 THB/user/month สำหรับ Pro
  
- **Reallocate Opportunities**: ระหว่าง departments
  - จำนวน: ~50-100+ opportunities
  - เงินประหยัด: ~1500 THB/user
  - แสดง pending requests
  
- **Downgrade Opportunities**: Users ที่ใช้ tier แพงแต่ไม่จำเป็น
  - จำนวน: ~100-200+ opportunities
  - เงินประหยัด: ~2000-2500 THB สำหรับ Enterprise
  - เงินประหยัด: ~800-1000 THB สำหรับ Pro

### ผลลัพธ์โดยรวม
- **Total Recommendations**: ~400-600 recommendations
- **Total Potential Savings**: ~500,000 - 1,000,000+ THB/month
PaymentMatrix
- **Actions**: revoke, reallocate, downgrade
- **Priority**: 1-8 (สูงถึงต่ำ)

## ราคา Licenses ที่ตั้งไว้

| Tier | ราคา (THB/month) | Target Category |
|------|------------------|-----------------|
| Basic | 800 | Productivity, Collaboration |
| Pro | 2,000 | Productivity, Collaboration, DevOps |
| Enterprise | 3,500 | Productivity, Collaboration, DevOps, Security |

## การตรวจสอบผลลัพธ์

หลังจากรัน script แล้ว สามารถตรวจสอบได้ด้วย SQL:

```sql
-- ดูจำนวน recommendations ทั้งหมด
SELECT 
  action, 
  COUNT(*) as count,
  SUM(impact_saving_amt) as total_savings
FROM recommendations
WHERE type = 'seat_opt'
GROUP BY action;

-- ดู details ของ recommendations
SELECT 
  action,
  priority,
  impact_saving_amt,
  reason_json->>'reason' as reason,
  reason_json->>'app_name' as app_name
FROM recommendations
WHERE type = 'seat_opt'
ORDER BY priority DESC, impact_saving_amt DESC
LIMIT 20;
```

## Troubleshooting

### ✅ แก้ไขแล้ว: "duplicate key value violates unique constraint"
- Scripts ได้แก้ไขแล้ว รองรับการรันซ้ำได้โดยอัตโนมัติ
- จะลบข้อมูลเก่าก่อน insert ใหม่ทุกครั้งที่รัน

### ข้อมูลน้อยเกินไป
- เพิ่มจำนวน generate_series (จาก 8 เป็น 15-20)
- ลดค่า random() threshold (จาก 0.7 เป็น 0.5)

### ข้อมูลมากเกินไป
- ลดจำนวน generate_series
- เพิ่มค่า random() threshold

### ถ้าต้องการ reset ทั้งหมด
```sql
-- ลบข้อมูลทั้งหมดที่เกี่ยวข้อง
DELETE FROM recommendations WHERE type = 'seat_opt';
DELETE FROM license_assignments WHERE reason IN ('seat_opt_mock_inactive', 'massive_seat_opt_revoke');
DELETE FROM users WHERE username LIKE 'inactive_%' OR username LIKE 'user_inactive_%';
```
