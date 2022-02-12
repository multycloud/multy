# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

"""
Purpose

Shows how to implement an AWS Lambda function that handles input from direct
invocation.
"""

import logging
import math
import json

logger = logging.getLogger()
logger.setLevel(logging.INFO)

# Define a list of Python lambda functions that are called by this AWS Lambda function.
ACTIONS = {
    'square': lambda x: x * x,
    'square root': lambda x: math.sqrt(x),
    'increment': lambda x: x + 1,
    'decrement': lambda x: x - 1,
}


def lambda_handler(event, context):
    """
    Accepts an action and a number, performs the specified action on the number,
    and returns the result.

    :param event: The event dict that contains the parameters sent when the function
                  is invoked.
    :param context: The context in which the function is called.
    :return: The result of the specified action.
    """
    logger.info('Event: %s', event)

    logging.info('Python HTTP trigger function processed a request.')

    if event['path'] != "/":
        return {
            "statusCode": 404,
            "headers": {
                "Content-Type": "application/json",
            },
            "body": json.dumps(event, indent=4),
        }

    name = None
    if event['queryStringParameters']:
        name = event['queryStringParameters']['name']

    logging.info('Processed name.')

    if name:
        return {
                "statusCode": 200,
            "body" : f"Hello, {name}. This HTTP triggered function executed successfully."
            }
    else:
        return {
            "statusCode": 200,
            "body" : "This HTTP triggered function executed successfully. Pass a name in the query string or in the request body for a personalized response."
        }
