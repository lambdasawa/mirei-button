import * as cdk from "@aws-cdk/core";
import { Bucket } from "@aws-cdk/aws-s3";
import { CloudFrontWebDistribution } from "@aws-cdk/aws-cloudfront";
import { CfnOutput } from "@aws-cdk/core";

export class MireiButtonStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const bucket = new Bucket(this, "Public", {});

    const distribution = new CloudFrontWebDistribution(
      this,
      "PublicDistribution",
      {
        originConfigs: [
          {
            s3OriginSource: {
              s3BucketSource: bucket
            },
            behaviors: [{ isDefaultBehavior: true }]
          }
        ]
      }
    );

    new CfnOutput(this, "BucketName", {
      exportName: "BucketName",
      value: bucket.bucketName
    });
    new CfnOutput(this, "PublicURL", {
      exportName: "PublicURL",
      value: `https://${distribution.domainName}`
    });
  }
}
