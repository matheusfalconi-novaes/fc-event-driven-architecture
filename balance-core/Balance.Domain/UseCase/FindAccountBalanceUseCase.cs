using Balance.Domain.Gateway;

namespace Balance.Domain.UseCase;

public record FindAccountBalanceInputDto(string AccountId);

public record FindAccountBalanceOutputDto(float Balance);

public class FindAccountBalanceUseCase(IBalanceGateway balanceGateway)
{
    public async Task<FindAccountBalanceOutputDto?> Execute(FindAccountBalanceInputDto dto)
    {
        var balance = await balanceGateway.FindBalanceByAccountId(dto.AccountId);
        return balance != null ? new FindAccountBalanceOutputDto(balance.Balance) : null;
    }
}