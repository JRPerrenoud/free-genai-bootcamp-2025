## Functional Requirements
Company wants to keep infrastruture on premise. They have already invested the capital for the equipment and do not want to pay more for a different infrastructure. They are also concerned about the privacy of their data and want to keep that localized to within company walls.  

This will be running on single server Xeon processor with Guadi 3 AI HW Accelerator

Customer is a secondary language learning high school located in Utah and will support 250-400 active students in total, but not all at the same time. 


## Assumptions
Open-source LLMs should be powerful enough to run on their HW.

Budget allocated to purchase copyrighted materiel **should** be enough

## Data Strategy
Copyrighted materials will need to be purchased

## Considerations
First pick right now with out much more thought is using IBM Granite because it is truley open-source with traceable training data (may help with copyright costs/issues) -- but willing to look at other models and data sets with same constraints.  (i.e. Meta Llama 8B)

https://huggingface.co/ibm-granite

https://huggingface.co/meta-llama/Meta-Llama-3-8B

## Use Cases
Students will work on language learning activies on terminals in a lab that access the lanague learning program on VMs running Windows.  They will be able to access the app thru Edge browswer.

