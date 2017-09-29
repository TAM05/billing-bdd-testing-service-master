# billing-bdd-testing-service

Runs the scenarios defined in feature files under `features` directory.

WIP!

## what is this?

Spawns appropriate telco usage-events and puts them to the entry-point of the rating system.
Waits for a configured duration then checks whether the event was rated correctly or not.

To trigger a run hit the `/api/1.0/run-test` endpoint with a `POST` request.

Account used for testing is:

```bson
{
	"accountNumber": "bd766d3e-3e92-4d4b-b837-8372343eee5e_TestAccount",
	"dob": null,
	"name": {
		"fullname": "Test Account",
		"title": "",
		"initials": "",
		"firstname": "",
		"surname": ""
	},
	"endDate": null,
	"liveDate": ISODate("2017-09-21T00:00:00Z"),
	"address": null,
	"rating": "good",
	"isOnPaymentPlan": true,
	"contractType": "residential",
	"services": [{
		"serviceId": "cc55ea4c-563d-4194-a17e-4e1e21199d30_Test",
		"label": "+440000000000_Test_030S_1",
		"endDate": null,
		"liveDate": ISODate("2017-09-21T00:00:00Z"),
		"tariff": "030S",
		"subTariff": "111111",
	}, {
		"serviceId": "2d6dd62e-86fd-4c07-86e3-63d17905d456_Test",
		"label": "+440000000000_Test_030S_2",
		"endDate": null,
		"liveDate": ISODate("2007-09-21T00:00:00Z"),
		"tariff": "030S",
		"subTariff": "222222"
	}]
}}
```
