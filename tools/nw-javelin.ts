import { AwsClient } from 'aws4fetch'
import { program } from 'commander'
import dotenv from 'dotenv'
import fs from 'fs'
import path from 'path'
import SteamUser from 'steam-user'

dotenv.config()

program
  .command('update')
  .description('Updates javelin config')
  .requiredOption('-u, --user <user>', 'Account user name', process.env.STEAM_USER)
  .requiredOption('-p, --pass <pass>', 'Account password', process.env.STEAM_PASS)
  .requiredOption('-o, --out <out>', 'Output directory', './tmp/javelin')
  .action(async ({ user, pass, out }: { user: string; pass: string; out: string }) => {
    fs.mkdirSync(out, { recursive: true })
    await resolveNwSession(user, pass)
      .then((session) => {
        console.log('login OK')
        return fetchCredentials(session.accessToken)
      })
      .then((credentials) => {
        console.log('credentials OK')
        return downloadConfigFiles(credentials, out)
      })
      .catch((error) => {
        console.error(error)
      })
  })

program.parse(process.argv)

async function resolveNwSession(user: string, pass: string) {
  return new Promise<NwSession>((resolve, reject) => {
    const client = new SteamUser()
    client.once('error', (error) => {
      reject(error)
      client.logOff()
    })
    client.once('loggedOn', async () => {
      await client
        .createAuthSessionTicket(1063730)
        .then(({ sessionTicket }) => {
          return createNwSession(sessionTicket.toString('hex'))
        })
        .then(resolve)
        .catch(reject)
      client.logOff()
    })
    client.logOn({
      accountName: user,
      password: pass,
    })
  })
}

type NwSession = {
  accessToken: string
  accessTokenExpirationDate: number
  expiresIn: number
  fallbackToken: string
  platformAccount: {
    identityType: 'steam'
    identityId: string
    personaId: string
    ageGroup: string
    platform: 'steam'
  }
  account: Record<string, any>
  auxiliaryPlatformAccounts: []
  penalties: []
  region: 'DE'
}

async function createNwSession(authToken: string): Promise<NwSession> {
  return fetch(`https://tokenservice.amazongames.com/games/new-world/tokens`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'User-Agent': 'OmniSDK/1.6.5/Windows',
      'Accept-Language': 'en-US',
    },
    body: JSON.stringify({
      platformType: 'steam',
      platformAuth: {
        ticket: authToken,
      },
      attributionData: {
        os: 'Windows 11',
        resolution: '1920x1080',
        language: 'en-US',
        platform: 'Windows',
        created: new Date().getTime() / 1000,
      },
    }),
  }).then((r) => r.json())
}

type OmniCredentials = {
  accessKeyId: string
  secretAccessKey: string
  sessionToken: string
  expiration: number
}

async function fetchCredentials(nwToken: string): Promise<OmniCredentials> {
  return fetch(`https://d2oeuvxi3kfsrw.cloudfront.net/prod/credentials/omni`, {
    headers: {
      'User-Agent': 'aws-sdk-cpp/1.7.193 Windows/10.0.22621.2506 AMD64 MSVC/1929',
      Authorization: 'Bearer ' + nwToken,
    },
  }).then((r) => r.json())
}

const JAVELIN_BASE = 'https://ags-javelin-remote-config.s3.amazonaws.com/'
const JAVELIN_URIS = [
  'applications/public/configuration-sets/CognitoId/us-east-1:0de93f7c-814e-4d48-a8fd-b6efc2f9ed93/versionless',
  'applications/public/configuration-sets/CognitoId/us-east-1:6ddd6498-c7d2-4df5-9ac6-855f140a3f09/versionless',
  'applications/public/configuration-sets/CognitoId/us-east-1:c654c4c8-bb6a-4e35-8200-b02bdd46b897/versionless',
  'applications/public/configuration-sets/ProductId/STEAM_APP_ID.1063730/versionless',
  'applications/public/configuration-sets/ProductId/STEAM_APP_ID.1205550/versionless',
  'applications/public/configuration-sets/RegionId/fra-prod/versionless',
  'applications/public/configuration-sets/RegionId/iad-prod/versionless',
  'applications/publicGameplay/configuration-sets/CognitoId/us-east-1:0de93f7c-814e-4d48-a8fd-b6efc2f9ed93/versionless',
  'applications/publicGameplay/configuration-sets/CognitoId/us-east-1:6ddd6498-c7d2-4df5-9ac6-855f140a3f09/versionless',
  'applications/publicGameplay/configuration-sets/CognitoId/us-east-1:c654c4c8-bb6a-4e35-8200-b02bdd46b897/versionless',
  'applications/publicGameplay/configuration-sets/ProductId/STEAM_APP_ID.1063730/versionless',
  'applications/publicGameplay/configuration-sets/RegionId/fra-prod/versionless',
  'applications/publicGameplay/configuration-sets/RegionId/iad-prod/versionless',
  'applications/publicGameplay/configuration-sets/WorldId/960d6d52-390b-4592-bb8b-79548f0956a1/versionless',
]
const MARKETING_BASE = 'https://ags-nw-cms.s3.us-west-2.amazonaws.com/'
const MARKETING_URIS = ['marketingtiles/STEAM_APP_ID.1205550/metadata.json']

async function downloadConfigFiles(credentials: OmniCredentials, outDir: string) {
  const client = new AwsClient({
    accessKeyId: credentials.accessKeyId,
    secretAccessKey: credentials.secretAccessKey,
    sessionToken: credentials.sessionToken,
  })
  for (const uri of [...MARKETING_URIS, ...JAVELIN_URIS]) {
    const outFile = path.join(outDir, uri.replaceAll(/[\/.:]/gi, '_').toLowerCase() + '.json')
    const baseUrl = MARKETING_URIS.includes(uri) ? MARKETING_BASE : JAVELIN_BASE
    const url = baseUrl + uri
    console.log(uri)
    await client
      .fetch(url)
      .then((res) => res.text())
      .then(async (data) => fs.promises.writeFile(outFile, data))
      .catch((e) => {
        console.error(e)
      })
  }
}
