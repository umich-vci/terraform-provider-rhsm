resource "rhsm_allocation_entitlement" "entitlement" {
  allocation_uuid = "123e4567-e89b-12d3-a456-426655440000"
  pool            = "123e4567e89b12d3a456426655440000"
  quantity        = 5
}
