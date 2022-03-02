import logging

import azure.functions as func


def main(req: func.HttpRequest) -> func.HttpResponse:
    logging.info('Python HTTP trigger function processed a request.')

    name = req.params.get('name')
    if not name:
        try:
            req_body = req.get_json()
        except ValueError:
            pass
        else:
            name = req_body.get('name')

    if name:
        return func.HttpResponse(f"Hello, {name}. This HTTP triggered azure",
                                 status_code=200,
                                 headers={
                                     "Access-Control-Allow-Origin": "*",
                                 })
    else:
        return func.HttpResponse(
            "This HTTP triggered azure.",
            status_code=200,
            headers={
                "Access-Control-Allow-Origin": "*",
            }
        )
