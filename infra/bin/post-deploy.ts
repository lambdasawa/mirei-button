import * as aws from "aws-sdk";
import * as fs from "fs";
import * as path from "path";
import * as mimeTypes from "mime-types";

const findStack = (): Promise<aws.CloudFormation.Stack> => {
  return new aws.CloudFormation()
    .describeStacks({
      StackName: process.env.MB_STACK_NAME || ""
    })
    .promise()
    .then(stacks => (stacks.Stacks || [])[0]);
};

const findBucketName = (stack: aws.CloudFormation.Stack): string => {
  return (
    stack.Outputs?.find(o => o.OutputKey === "BucketName")?.OutputValue || ""
  );
};

const findDistributionID = (stack: aws.CloudFormation.Stack): string => {
  return (
    stack.Outputs?.find(o => o.OutputKey === "DistributionID")?.OutputValue ||
    ""
  );
};

const putObject = (
  bucketName: string,
  prefix: string,
  key: string
): Promise<aws.S3.PutObjectOutput> => {
  console.log({ msg: "putObject", key });

  const file = fs.readFileSync(key);

  const contentType = mimeTypes.lookup(key) || undefined;

  return new aws.S3()
    .putObject({
      Bucket: bucketName,
      Key: key.replace(prefix, ""),
      Body: file,
      ACL: "public-read",
      ContentType: contentType
    })
    .promise();
};

const walk = (
  dir: string,
  callback: (key: string) => Promise<unknown>
): Promise<unknown> => {
  return Promise.all(
    fs.readdirSync(dir).map(async name => {
      const filePath = path.join(dir, name);

      const stat = fs.statSync(filePath);

      if (stat.isFile()) {
        return callback(filePath);
      }
      if (stat.isDirectory()) {
        return walk(filePath, callback);
      }
    })
  );
};

const invalidateCloudFront = (distributionID: string): Promise<void> => {
  const cf = new aws.CloudFront();

  const paths = [
    "/js/*",
    "/index.html",
    "/metadata.json",
    "/css/*",
    "/favicon.ico",
    "/img/*"
  ];
  return cf
    .createInvalidation({
      DistributionId: distributionID,
      InvalidationBatch: {
        Paths: {
          Quantity: paths.length,
          Items: paths
        },
        CallerReference: String(new Date().getTime())
      }
    })
    .promise()
    .then(() => undefined);
};

const main = async (): Promise<unknown> => {
  aws.config.logger = console;

  const stack = await findStack();
  const bucketName = findBucketName(stack);

  const prefix = "../public/dist";
  await walk(prefix, async key => putObject(bucketName, prefix, key));

  const distributionID = findDistributionID(stack);
  await invalidateCloudFront(distributionID);

  return;
};

main()
  .then(() => {
    process.exit(0);
  })
  .catch(e => {
    console.error(e);
    process.exit(1);
  });
