use etcd_client::{Client, Error};

#[tokio::main]
async fn main() -> Result<(), Error> {
    let mut client = Client::connect(["localhost:2379"], None).await?;
    // put kv
    client.put("/web/foo", "bar", None).await?;
    // get kv
    let resp = client.get("/web/foo", None).await?;
    if let Some(kv) = resp.kvs().first() {
        println!("Get kv: {{{}: {}}}", kv.key_str()?, kv.value_str()?);
    }

    Ok(())
}