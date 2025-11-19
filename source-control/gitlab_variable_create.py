import requests

GITLAB_URL = "https://gitlab.com/api/v4"
TOKEN = "your_token_here"

headers = {
    "PRIVATE-TOKEN": TOKEN,
    "Content-Type": "application/json"
}

# List of project IDs
project_ids = [1234, 5678, 91011]

# Variables to create in each project
variables = {
    "DB_HOST": "10.0.0.1",
    "DB_USER": "admin",
    "DB_PASS": "password123",
    "ENVIRONMENT": "production"
}

for project_id in project_ids:
    print(f"\n🔧 Creating variables for Project {project_id}")

    for key, value in variables.items():
        payload = {
            "key": key,
            "value": value,
            "masked": False,
            "protected": False,
            "environment_scope": "*"
        }

        response = requests.post(
            f"{GITLAB_URL}/projects/{project_id}/variables",
            headers=headers,
            json=payload
        )

        if response.status_code == 201:
            print(f"   ✔ Created {key}")
        elif response.status_code == 400 and "Key has already been taken" in response.text:
            print(f"   ⚠ {key} already exists, skipping")
        else:
            print(f"   ❌ Error creating {key}: {response.text}")

