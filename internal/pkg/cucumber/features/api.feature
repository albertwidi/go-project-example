Feature: get book detail
    In order eto get book
    As an API user
    I need to be able to request book detail

    Scenario: get book detail
        Given set request header "Content-Type: application/json"
        And set request body:
            """
            {
                "book_id": 10
            }
            """
        When I send "POST" request to "http://127.0.0.1:9863/v1/book/detail"
        Then the response code should be 200
        And the response header "Content-Type" should be "text/plain; charset=utf-8"
        And the response should match json:
            """
            {
                "book_id": 10,
                "name": "testing",
                "author": "what?"
            }
            """

    Scenario: get book detail with endpoint mapping and path
        Given set request header "Content-Type: application/json"
        And set request body:
            """
            {
                "book_id": 10
            }
            """
        When I send "POST" request to "book service" with path "/v1/book/detail"
        Then the response code should be 200
        And the response header "Content-Type" should be "text/plain; charset=utf-8"
        And the response should match json:
            """
            {
                "book_id": 10,
                "name": "testing",
                "author": "what?"
            }
            """
    
    Scenario Outline: this is an example of scenario outlines
        Given set request header "<request_content_type>"
        And set request body:
            """
            <request_body>
            """
        When I send "<method>" request to "<service> service" with path "<path>"
        Then the response code should be <response_code>
        And the response header "<response_header_key>" should be "<response_header_value>"
        And the response should match json:
            """
            <response_body>
            """
        Examples:
            | request_content_type           | request_body     | method | service | path            | response_code | response_header_key | response_header_value     | response_body                                         |
            | Content-Type: application/json | {"book_id": 10}  | POST   | book    | /v1/book/detail | 200           | Content-Type        | text/plain; charset=utf-8 | {"book_id": 10, "name": "testing", "author": "what?"} |
            | Content-Type: application/json | {"book_id": 20}  | POST   | book    | /v1/book/detail | 200           | Content-Type        | text/plain; charset=utf-8 | {"book_id": 20, "name": "testing", "author": "what?"} |