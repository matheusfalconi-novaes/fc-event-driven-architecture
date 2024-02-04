using Balance.Domain.Gateway;
using InfluxDB.Client;
using InfluxDB.Client.Api.Domain;
using InfluxDB.Client.Writes;
using Microsoft.Extensions.Options;

namespace Balance.Infrastructure.Gateway;

public record BalanceInfluxDbConfiguration
{
    public string BucketName { get; set; } = null!;
    public string Organization { get; set; } = null!;
    public string Url { get; set; } = null!;
    public string Token { get; set; } = null!;
}

public class BalanceInfluxDb(IOptions<BalanceInfluxDbConfiguration> configuration) : IBalanceGateway
{
    private readonly InfluxDBClient _influxDbClient = new(configuration.Value.Url, configuration.Value.Token);
    private readonly string _bucketName = configuration.Value.BucketName;
    private readonly string _organization = configuration.Value.Organization;
    private const string Measurement = "balances";

    public Task SaveBalances(IEnumerable<BalanceData> balancesData, CancellationToken cancellationToken = default)
    {
        var write = _influxDbClient.GetWriteApiAsync();
        var measurement = PointData.Measurement(Measurement);
        var points = balancesData
            .Select(balance => measurement
                .Tag("account_id", balance.AccountId)
                .Field("balance_account", balance.Balance)
                .Timestamp(balance.OperationDateTime, WritePrecision.Ns))
            .ToArray();
        return write.WritePointsAsync(points, _bucketName, _organization, cancellationToken);
    }

    public async Task<BalanceData?> FindBalanceByAccountId(string accountId)
    {
        var query = $"from(bucket: \"{_bucketName}\") " +
                    $"|> range(start: 0) " +
                    $"|> filter (fn: (r) => " +
                    $"r._measurement == \"{Measurement}\" and " +
                    $"r.account_id == \"{accountId}\")";

        var lastBalance = (await _influxDbClient.GetQueryApi().QueryAsync(query, _organization))
            ?.SelectMany(table => table.Records)
            .MaxBy(record => record.GetTimeInDateTime());

        if (!float.TryParse(lastBalance?.GetValue().ToString(), out var parsedValue) ||
            !DateTime.TryParse(lastBalance.GetTimeInDateTime().ToString(), out var parsedDate))
            return null;
        
        return new BalanceData(accountId, parsedValue, parsedDate);
    }
}