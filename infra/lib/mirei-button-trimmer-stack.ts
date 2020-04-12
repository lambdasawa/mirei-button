import * as path from "path";
import * as cdk from "@aws-cdk/core";
import * as ec2 from "@aws-cdk/aws-ec2";
import * as ecr from "@aws-cdk/aws-ecr";
import * as ecsAssets from "@aws-cdk/aws-ecr-assets";
import * as ecs from "@aws-cdk/aws-ecs";
import * as ecsPatterns from "@aws-cdk/aws-ecs-patterns";
import * as route53 from "@aws-cdk/aws-route53";
import * as route53Targets from "@aws-cdk/aws-route53-targets";

export class MireiButtonTrimmerStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const asset = new ecsAssets.DockerImageAsset(this, "DockerImageAsset", {
      directory: path.join(__dirname, "../../trimmer")
    });

    const vpc = ec2.Vpc.fromLookup(this, "VPC", {
      isDefault: true
    });

    const cluster = new ecs.Cluster(this, "Cluster", {
      vpc: vpc
    });

    const loadBalancedFargateService = new ecsPatterns.ApplicationLoadBalancedFargateService(
      this,
      "Service",
      {
        cluster,
        memoryLimitMiB: 4096,
        cpu: 2048,
        assignPublicIp: true,
        taskImageOptions: {
          image: ecs.ContainerImage.fromEcrRepository(
            asset.repository,
            asset.sourceHash
          ),
          containerPort: 3011,
          environment: this.getEnvironment()
        }
      }
    );

    const domainName = "mirei-button.net";
    const recordName = ["trimmer", "prod"].join("."); // TODO: switch stage

    const hostedZone = route53.HostedZone.fromLookup(this, "HostedZone", {
      domainName
    });

    new route53.ARecord(this, "AliasRecord", {
      zone: hostedZone,
      recordName,
      target: route53.RecordTarget.fromAlias(
        new route53Targets.LoadBalancerTarget(
          loadBalancedFargateService.loadBalancer
        )
      )
    });

    new cdk.CfnOutput(this, "TrimmerPublicURL", {
      exportName: "TrimmerPublicURL",
      value: `http://${recordName}.${domainName}`
    });
  }

  getEnvironment(): Record<string, string> {
    const environment: Record<string, string> = {};
    const keys = [
      "MB_STAGE",
      "MB_STACK_NAME",
      "MB_YOUTUBEDL_BIN_PATH",
      "MB_FFMPEG_BIN_PATH",
      "MB_SOX_BIN_PATH",
      "MB_SESSION_SECRET",
      "MB_TWITTER_CONSUMER_KEY",
      "MB_TWITTER_CONSUMER_SECRET",
      "MB_TWTITER_CALLBACK",
      "MB_TWITTER_SCREENNAME",
      "MB_BUCKET",
      "MB_DISTRIBUTION_ID"
    ];

    keys.forEach(key => {
      environment[key] = process.env[key] || "";
    });

    return environment;
  }
}
