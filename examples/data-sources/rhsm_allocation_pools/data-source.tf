data "rhsm_allocation_pools" "pools" {
  allocation_uuid = "123e4567-e89b-12d3-a456-426655440000"
  future          = true
}
