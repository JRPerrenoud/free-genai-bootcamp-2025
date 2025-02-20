import boto3
import json

def test_aws_access():
    print("Testing AWS Access...")
    
    # Test STS access
    try:
        sts = boto3.client('sts')
        identity = sts.get_caller_identity()
        print("\nAWS Identity:")
        print(f"Account: {identity['Account']}")
        print(f"User/Role ARN: {identity['Arn']}")
    except Exception as e:
        print(f"\nError accessing STS: {str(e)}")
    
    # Test Bedrock access
    try:
        bedrock = boto3.client('bedrock-runtime', region_name='us-east-2')
        print("\nBedrock client created successfully")
        
        # Try to invoke the model
        body = json.dumps({"inputText": "test"})
        response = bedrock.invoke_model(
            modelId="amazon.titan-embed-text-v2:0",
            body=body
        )
        print("\nBedrock model invoked successfully!")
    except Exception as e:
        print(f"\nError accessing Bedrock: {str(e)}")

if __name__ == "__main__":
    test_aws_access()
