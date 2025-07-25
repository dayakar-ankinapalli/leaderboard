import sys
from pathlib import Path
from typing import Optional

from pydantic import (
    BaseModel,
    Field,
    PositiveInt,
    SecretStr,
    ValidationError,
)
from ruamel.yaml import YAML
from ruamel.yaml.error import YAMLError

# 1. Define the Schema for Validation (The "Accuracy" part)
# Using Pydantic models to define the expected structure and types of our config.
# This is far more robust and accurate than just loading into a dictionary.

class DatabaseConfig(BaseModel):
    """Defines the schema for the database connection details."""
    host: str
    port: int = Field(gt=0, lt=65536)  # Port must be in the valid range
    user: str
    password: SecretStr  # Pydantic's type for secrets, prevents them from being logged

class FeatureConfig(BaseModel):
    """Defines the schema for the feature flags."""
    enable_beta_feature: bool
    enable_dark_mode: bool
    max_users: PositiveInt  # Ensures the value is an integer > 0

class AppConfig(BaseModel):
    """The top-level configuration model that aggregates all other parts."""
    database: DatabaseConfig
    api_key: SecretStr
    features: FeatureConfig


# 2. Create an Efficient and Safe Loader Function

def load_and_validate_config(config_path: Path) -> Optional[AppConfig]:
    """
    Loads a YAML configuration file efficiently and validates it against the AppConfig schema.

    Args:
        config_path: The path to the YAML configuration file.

    Returns:
        An AppConfig object if loading and validation are successful, otherwise None.
    """
    if not config_path.is_file():
        print(f"Error: Configuration file not found at {config_path}", file=sys.stderr)
        return None

    # Use ruamel.yaml with the 'safe' loader. It automatically uses the
    # fast C-based libyaml if available, making it highly efficient.
    yaml = YAML(typ='safe')

    try:
        print(f"Loading configuration from: {config_path}")
        data = yaml.load(config_path)

        print("Configuration loaded, now validating for accuracy...")
        config = AppConfig.model_validate(data)
        print("Configuration is valid!")
        return config

    except YAMLError as e:
        print(f"Error parsing YAML file: {e}", file=sys.stderr)
    except ValidationError as e:
        print(f"Configuration validation error:\n{e}", file=sys.stderr)
    except Exception as e:
        print(f"An unexpected error occurred: {e}", file=sys.stderr)
    
    return None


# 3. Example Usage
if __name__ == "__main__":
    config_file = Path(__file__).parent / "config.yaml"
    app_config = load_and_validate_config(config_file)

    if app_config:
        print("\n--- Accessing Validated Configuration ---")
        # Accessing data is now type-safe and your IDE provides autocompletion.
        print(f"Database Host: {app_config.database.host}")
        print(f"Max Users Feature: {app_config.features.max_users}")

        # Pydantic's SecretStr hides the value when printed for security
        print(f"API Key: {app_config.api_key}")
        # To get the actual value, you must explicitly call get_secret_value()
        print(f"Actual API Key Value: '{app_config.api_key.get_secret_value()}'")

        if app_config.features.enable_beta_feature:
            print("Beta feature is enabled.")