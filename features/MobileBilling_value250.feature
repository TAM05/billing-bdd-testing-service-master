#
# Feature: Mobile Billing_value250
#
# Created with BDD Editor on: 22 August, 2017
#
# Please follow us at @bddeditor AND if you find this tool useful please share with friends and colleagues!
#

Feature: Mobile Billing_value250
  As A bla
  I want blah
  So That bla-blah

  @Value250
  Scenario Outline: Value250_MO to MT
    Given MO is on tariff *Value250
    And MO with number <CLI> connects to MT via <Retail band> and <Communication type>
    Then get result and <Total rate>

    Examples:
      |       MO         |              CLI             |  Retail band  |    MT    | Communication type   | Result  | Price | Total charge | Call set up fee | Total rate |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     m121      | On net   | voice (300 seconds)  | success | .30p  |     1.50     |      0.05       |    1.55    |

  # used for application testing
  @test
  Scenario Outline: Value250_MO to MT
    Given MO is on tariff *Value250
    And MO with number <CLI> connects to MT via <Retail band> and <Communication type>
    Then get result and <Total rate>

    Examples:
      |       MO         |              CLI             |  Retail band  |    MT    | Communication type   | Result  | Price | Total charge | Call set up fee | Total rate |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     m121      | On net   | voice (300 seconds)  | success | .25p  |     1.25     |      0.05       |    1.30    |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     smson     | On net   | sms (1 unit)         | success | .00p  |     0.00     |      0.00       |    0.00    |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     pic       | On net   | mms (1 unit)         | success | .25p  |     0.25     |      0.05       |    0.30    |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     vm        | On net   | voice (1 unit)       | success | .10p  |     0.10     |      0.05       |    0.15    |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     gprs      | On net   | data (10000 kb)      | success | 0.04p |     0.40     |      0.05       |    0.45    |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     mobil     | Off net  | voice (360 seconds)  | success | .25p  |     1.50     |      0.05       |    1.55    |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     smsof     | Off net  | sms (1 unit)         | success | .25p  |     0.00     |      0.05       |    0.00    |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     pic       | Off net  | mms (1 unit)         | success | .25p  |     0.25     |      0.05       |    0.30    |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     vm        | Off net  | voice (1 unit)       | success | .10p  |     0.10     |      0.05       |    0.15    |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     roamt     | Off net  | sms (1 unit)         | success | 0.25p |     0.25     |      0.05       |    0.30    |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     zpic      | Off net  | mms (1 unit)         | success | 0.37p |     0.37     |      0.05       |    0.42    |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     zn08i     | Off net  | voice(1 unit)        | success | 1.75p |     1.75     |      0.05       |    1.75    |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     rgp05     | Off net  | data (5000 kb)       | success | 7.00p |     35.00    |      0.05       |    0.45    |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     zn07i     | Off net  | voice (120 seconds)  | success | 2.10p |     4.20     |      0.00       |    4.20    |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     intt      | Off net  | sms (1 unit)         | success | 0.20p |     0.20     |      0.05       |    0.25    |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     zpic      | Off net  | mms (1 unit)         | success | .37p  |     0.00     |      0.05       |    0.42    |
      | Pay monthly (UW) |   +440000000000_Test_030S_1  |     vm        | Off net  | voice (1 unit)       | success | .12p  |     0.12     |      0.05       |    0.17    |
        
      