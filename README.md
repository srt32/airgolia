# airgolia

Read data from Airtable and post it to an Algolia index.

## Example usage

`AIRTABLE_API_KEY=foo AIRTABLE_TABLE_NAME=bar AIRTABLE_WORKSPACE_ID=foo ALGOLIA_APP_ID=foo ALGOLIA_API_KEY=foo ALGOLIA_INDEX_NAME=foo run main.go`

## Notes

* When run, this task will fetch all records in the Airtable table and create or update (based on id).
* All fields are synced.
* The Airtable calls do not paginate so if there are too many records they won't sync (for now).
