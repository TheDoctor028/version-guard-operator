from flask import Flask, request

app = Flask("test-api")


@app.route('/', defaults={'path': ''}, methods=['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS'])
@app.route('/<path:path>', methods=['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS'])
def catch_all(path):
    # Log the incoming request method and data (if any)
    method = request.method
    data = request.data

    # Print the request method and data to the console
    print(f'Received {method} request for path: /{path}')
    print(f'Data: {data.decode()}')

    # You can process the request data or return a custom response here
    # For simplicity, we are just printing the request data

    return f'Received {method} request for path: /{path}\nData: {data.decode()}'


if __name__ == '__main__':
    app.run(debug=True)
