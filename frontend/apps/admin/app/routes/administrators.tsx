import { UserAdmin } from "./users"

export function meta() {
  return [{ title: "管理账号 · Referral Admin" }]
}

export default function Administrators() {
  return <UserAdmin accountType="admin" />
}
