namespace Balance.Domain.Gateway;

public record BalanceData(string AccountId, float Balance, DateTime OperationDateTime);

public interface IBalanceGateway
{
    Task SaveBalances(IEnumerable<BalanceData> balancesData, CancellationToken cancellationToken = default);
    Task<BalanceData?> FindBalanceByAccountId(string accountId);
}